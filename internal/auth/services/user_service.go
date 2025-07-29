package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/fraineri/plexo_backend/internal/auth/persistance"
)

type UserService interface {
	CheckIfEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, firstName, lastName, email, hashedPassword string) (*entities.User, error)
	ActivateUser(ctx context.Context, userID string) error
	Login(ctx context.Context, email, password string) (*entities.User, error) // New method
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

func (s *userService) Login(ctx context.Context, email, password string) (*entities.User, error) {
	// 1. Find user by email
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Check user status
	if user.Status == "PENDING_VERIFICATION" {
		return nil, errors.New("account not verified")
	}
	if user.Status == "SUSPENDED" {
		return nil, errors.New("account suspended")
	}
	if user.Status == "BLOCKED" {
		if user.BlockedUntil != nil && time.Now().Unix() < *user.BlockedUntil {
			return nil, fmt.Errorf("account blocked until %s", time.Unix(*user.BlockedUntil, 0).Format(time.RFC1123))
		}
		// If block expired, unblock the user
		user.Status = "ACTIVE"
		user.BlockedUntil = nil
		if err := s.userRepository.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to unblock user: %w", err)
		}
	}

	// 3. Check password
	if !s.hashingService.CheckPasswordHash(password, user.PasswordHash) {
		// Password incorrect, record failed attempt and check for lockout
		_ = s.loginAttemptRepository.Create(ctx, &entities.LoginAttempt{UserID: user.ID, IsSuccessful: false})
		if err := s.handleFailedLogin(ctx, user); err != nil {
			return nil, err // Return the specific error (e.g., account now blocked)
		}
		return nil, errors.New("invalid credentials")
	}

	// 4. Successful login
	_ = s.loginAttemptRepository.Create(ctx, &entities.LoginAttempt{UserID: user.ID, IsSuccessful: true})
	return user, nil
}

func (s *userService) handleFailedLogin(ctx context.Context, user *entities.User) error {
	attempts, err := s.loginAttemptRepository.FindLastByUser(ctx, user.ID, 5)
	if err != nil {
		return fmt.Errorf("could not retrieve login attempts: %w", err)
	}

	if len(attempts) < 5 {
		return nil // Not enough attempts to trigger a lock
	}

	for _, attempt := range attempts {
		if attempt.IsSuccessful {
			return nil // A successful login breaks the chain
		}
	}

	// All 5 are failures, block the account
	blockDuration := 15 * time.Minute
	blockedUntil := time.Now().Add(blockDuration).Unix()
	user.Status = "BLOCKED"
	user.BlockedUntil = &blockedUntil

	if err := s.userRepository.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to lock user account: %w", err)
	}

	return fmt.Errorf("account has been blocked for %v due to too many failed login attempts", blockDuration)
}
