package services

import (
	"context"
	"fmt"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
)

type UserService interface {
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, firstName, lastName, email, hashedPassword string) (*entities.User, error)
	ActivateUser(ctx context.Context, userID string) error
	Authenticate(ctx context.Context, email, password string) (*entities.User, error)
	RecordSuccessfulLogin(ctx context.Context, userID string) error
	RecordFailedLogin(ctx context.Context, userID string) error
}

type userService struct {
	userRepository         persistance.UserRepository
	loginAttemptRepository persistance.LoginAttemptRepository
	hashingService         HashingService
}

func NewUserService(
	userRepository persistance.UserRepository,
	loginAttemptRepository persistance.LoginAttemptRepository,
	hashingService HashingService,
) UserService {
	return &userService{
		userRepository:         userRepository,
		loginAttemptRepository: loginAttemptRepository,
		hashingService:         hashingService,
	}
}

func (s *userService) CheckIfEmailExists(ctx context.Context, email string) (bool, error) {
	existingUser, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check for existing user: %w", err)
	}
	return existingUser != nil, nil
}

func (s *userService) CreateUser(ctx context.Context, firstName, lastName, email, hashedPassword string) (*entities.User, error) {
	newUser := &entities.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: hashedPassword,
		Status:       "PENDING_VERIFICATION",
	}

	if err := s.userRepository.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %w", err)
	}
	return newUser, nil
}

func (s *userService) ActivateUser(ctx context.Context, userID string) error {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("could not find user to activate: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user with id %s not found", userID)
	}

	if user.Status != "PENDING_VERIFICATION" {
		return fmt.Errorf("user is not pending verification, current status: %s", user.Status)
	}

	user.Status = "ACTIVE"
	if err := s.userRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user status to active: %w", err)
	}
	return nil
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check user status
	if user.Status == "PENDING_VERIFICATION" {
		return nil, ErrAccountNotVerified
	}
	if user.Status == "SUSPENDED" {
		return nil, ErrAccountSuspended
	}
	if user.Status == "BLOCKED" {
		if user.BlockedUntil != nil && time.Now().Unix() < *user.BlockedUntil {
			return nil, ErrAccountBlocked
		}
		// If block expired, unblock the user before proceeding
		user.Status = "ACTIVE"
		user.BlockedUntil = nil
		if err := s.userRepository.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to unblock user: %w", err)
		}
	}

	if !s.hashingService.CheckPasswordHash(password, user.PasswordHash) {
		return user, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) RecordSuccessfulLogin(ctx context.Context, userID string) error {
	return s.loginAttemptRepository.Create(ctx, &entities.LoginAttempt{UserID: userID, IsSuccessful: true})
}

func (s *userService) RecordFailedLogin(ctx context.Context, userID string) error {
	// 1. Record the failed attempt.
	err := s.loginAttemptRepository.Create(ctx, &entities.LoginAttempt{UserID: userID, IsSuccessful: false})
	if err != nil {
		return fmt.Errorf("failed to record login attempt: %w", err)
	}

	// 2. Check for account lockout condition.
	attempts, err := s.loginAttemptRepository.FindLastByUser(ctx, userID, 5)
	if err != nil {
		return fmt.Errorf("could not retrieve login attempts: %w", err)
	}

	if len(attempts) < 5 {
		return nil // Not enough attempts to trigger a lock.
	}

	isConsecutiveFailure := true
	for _, attempt := range attempts {
		if attempt.IsSuccessful {
			isConsecutiveFailure = false
			break
		}
	}

	if !isConsecutiveFailure {
		return nil // A successful login broke the chain.
	}

	// 3. All 5 are consecutive failures, block the account.
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil || user == nil {
		return fmt.Errorf("failed to find user to lock: %w", err)
	}

	blockDuration := 15 * time.Minute
	blockedUntil := time.Now().Add(blockDuration).Unix()
	user.Status = "BLOCKED"
	user.BlockedUntil = &blockedUntil

	if err := s.userRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to lock user account: %w", err)
	}

	// Return a specific business error indicating lockout.
	return ErrAccountLockout
}
