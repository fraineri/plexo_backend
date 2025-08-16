package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fraineri/plexo_backend/internal/auth/persistance"
	"github.com/fraineri/plexo_backend/internal/auth/services"
	"github.com/fraineri/plexo_backend/internal/auth/usecases"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
	"github.com/fraineri/plexo_backend/internal/core/settings"
	"github.com/gorilla/mux"
)

type authDependencies struct {
	registerUserUseCase           usecases.RegisterUserUseCase
	verifyAccountUseCase          usecases.VerifyAccountUseCase
	resendVerificationOtpUseCase  usecases.ResendVerificationOtpUseCase
	requestPasswordResetUseCase   usecases.RequestPasswordResetUseCase
	resetPasswordUseCase          usecases.ResetPasswordUseCase
	loginUserUseCase              usecases.LoginUserUseCase
	createRoleUseCase             usecases.CreateRoleUseCase
	createPermissionUseCase       usecases.CreatePermissionUseCase
	assignPermissionToRoleUseCase usecases.AssignPermissionToRoleUseCase
	assignRoleToUserUseCase       usecases.AssignRoleToUserUseCase
}

func newAuthDependencies(db *sql.DB, settings *settings.Settings) *authDependencies {
	// The UnitOfWork is created once per request.
	unitOfWork := uow.NewUnitOfWork(db)

	// Instantiate repositories with the request-scoped UnitOfWork.
	userRepository := persistance.NewUserRepository(unitOfWork)
	otpRepository := persistance.NewOtpRepository(unitOfWork)
	loginAttemptRepository := persistance.NewLoginAttemptRepository(unitOfWork)
	roleRepository := persistance.NewRoleRepository(unitOfWork)
	permissionRepository := persistance.NewPermissionRepository(unitOfWork)
	rolePermissionRepository := persistance.NewRolePermissionRepository(unitOfWork)
	userRoleRepository := persistance.NewUserRoleRepository(unitOfWork)

	// Instantiate services.
	hashingService := services.NewHashingService()
	notificationService := services.NewLogNotificationService() // Replace with a real service in production.
	authOtpService := services.NewAuthOtpService(otpRepository, settings.OTP)
	jwtService := services.NewJWTService(settings.JWT.SecretKey, settings.JWT.SessionTimeMins)
	userService := services.NewUserService(userRepository, loginAttemptRepository, hashingService)

	// Instantiate use cases, injecting the UoW and services.
	registerUserUseCase := usecases.NewRegisterUser(unitOfWork, userService, authOtpService, hashingService, notificationService)
	verifyAccountUseCase := usecases.NewVerifyAccount(unitOfWork, userService, authOtpService, userRepository)
	resendVerificationOtpUseCase := usecases.NewResendVerificationOtp(unitOfWork, userRepository, authOtpService, notificationService)
	requestPasswordResetUseCase := usecases.NewRequestPasswordReset(unitOfWork, userRepository, authOtpService, notificationService)
	resetPasswordUseCase := usecases.NewResetPassword(unitOfWork, userService, authOtpService, userRepository)
	loginUserUseCase := usecases.NewLoginUser(unitOfWork, userService, jwtService, userRoleRepository)
	createRoleUseCase := usecases.NewCreateRole(unitOfWork, roleRepository)
	createPermissionUseCase := usecases.NewCreatePermission(unitOfWork, permissionRepository)
	assignPermissionToRoleUseCase := usecases.NewAssignPermissionToRole(unitOfWork, rolePermissionRepository)
	assignRoleToUserUseCase := usecases.NewAssignRoleToUser(unitOfWork, userRoleRepository)

	return &authDependencies{
		registerUserUseCase:           registerUserUseCase,
		verifyAccountUseCase:          verifyAccountUseCase,
		resendVerificationOtpUseCase:  resendVerificationOtpUseCase,
		requestPasswordResetUseCase:   requestPasswordResetUseCase,
		resetPasswordUseCase:          resetPasswordUseCase,
		loginUserUseCase:              loginUserUseCase,
		createRoleUseCase:             createRoleUseCase,
		createPermissionUseCase:       createPermissionUseCase,
		assignPermissionToRoleUseCase: assignPermissionToRoleUseCase,
		assignRoleToUserUseCase:       assignRoleToUserUseCase,
	}
}

type AuthHandler struct {
	db       *sql.DB
	settings *settings.Settings
}

func NewAuthHandler(db *sql.DB, settings *settings.Settings) *AuthHandler {
	return &AuthHandler{
		db:       db,
		settings: settings,
	}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.registerUser).Methods("POST")
	router.HandleFunc("/auth/verify-account", h.verifyAccount).Methods("POST")
	router.HandleFunc("/auth/resend-verification-otp", h.resendVerificationOtp).Methods("POST")
	router.HandleFunc("/auth/request-password-reset", h.requestPasswordReset).Methods("POST")
	router.HandleFunc("/auth/reset-password", h.resetPassword).Methods("POST")
	router.HandleFunc("/auth/login", h.login).Methods("POST")
	router.HandleFunc("/auth/roles", h.createRole).Methods("POST")
	router.HandleFunc("/auth/permissions", h.createPermission).Methods("POST")
	router.HandleFunc("/auth/roles/permissions", h.assignPermissionToRole).Methods("POST")
	router.HandleFunc("/auth/users/roles", h.assignRoleToUser).Methods("POST")

}

func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var dto RegisterUserRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)

	input := usecases.RegisterUserInput{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
	}

	_, err := deps.registerUserUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful. Please check your email for the verification code."}); err != nil {
		log.Printf("Encoding response error: %v", err)
	}
}

func (h *AuthHandler) verifyAccount(w http.ResponseWriter, r *http.Request) {
	var dto VerifyAccountRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)

	input := usecases.VerifyAccountInput{
		Email:   dto.Email,
		OTPCode: dto.OTPCode,
	}

	err := deps.verifyAccountUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Verification error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Account verified successfully."}); err != nil {
		log.Printf("Encoding response error: %v", err)
	}
}

func (h *AuthHandler) resendVerificationOtp(w http.ResponseWriter, r *http.Request) {
	var dto ResendVerificationOtpRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)

	input := usecases.ResendVerificationOtpInput{
		Email: dto.Email,
	}

	if err := deps.resendVerificationOtpUseCase.Execute(r.Context(), input); err != nil {
		// Even if an error occurs (e.g., rate limit), we return a generic success message
		// to prevent leaking information. The error should be logged for monitoring.
		log.Printf("Resend OTP error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "If an account with that email exists, a new verification code has been sent."}); err != nil {
		log.Printf("Encoding response error: %v", err)
	}
}

func (h *AuthHandler) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var dto RequestPasswordResetDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)
	input := usecases.RequestPasswordResetInput{Email: dto.Email}

	if err := deps.requestPasswordResetUseCase.Execute(r.Context(), input); err != nil {
		log.Printf("Request password reset error: %v", err)
		// Return a generic error to the client, but log the specific one.
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "If an account with that email exists, a password reset code has been sent."}); err != nil {
		log.Printf("Encoding response error: %v", err)
	}
}

func (h *AuthHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var dto ResetPasswordRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)
	input := usecases.ResetPasswordInput{
		Email:       dto.Email,
		OTPCode:     dto.OTPCode,
		NewPassword: dto.NewPassword,
	}

	if err := deps.resetPasswordUseCase.Execute(r.Context(), input); err != nil {
		log.Printf("Reset password error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Password has been reset successfully."}); err != nil {
		log.Printf("Encoding response error: %v", err)
	}
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var dto LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)

	input := usecases.LoginUserInput{
		Email:    dto.Email,
		Password: dto.Password,
	}

	result, err := deps.loginUserUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Login error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusUnauthorized) // 401 for login failures
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponseDTO{Token: result.Token}); err != nil {
		log.Printf("Encoding response error: %v", err)
	}
}

func (h *AuthHandler) createRole(w http.ResponseWriter, r *http.Request) {
	var dto CreateRoleRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)
	input := usecases.CreateRoleInput{
		Name:        dto.Name,
		Description: dto.Description,
	}

	_, err := deps.createRoleUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Create role error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) createPermission(w http.ResponseWriter, r *http.Request) {
	var dto CreatePermissionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)
	input := usecases.CreatePermissionInput{
		Action:      dto.Action,
		Description: dto.Description,
	}

	_, err := deps.createPermissionUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Create permission error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) assignPermissionToRole(w http.ResponseWriter, r *http.Request) {
	var dto AssignPermissionToRoleRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)
	input := usecases.AssignPermissionToRoleInput{
		RoleID:       dto.RoleID,
		PermissionID: dto.PermissionID,
		Status:       dto.Status,
	}

	err := deps.assignPermissionToRoleUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Assign permission to role error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) assignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var dto AssignRoleToUserRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	deps := newAuthDependencies(h.db, h.settings)
	input := usecases.AssignRoleToUserInput{
		UserID: dto.UserID,
		RoleID: dto.RoleID,
	}

	err := deps.assignRoleToUserUseCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("Assign role to user error: %v", err)
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
