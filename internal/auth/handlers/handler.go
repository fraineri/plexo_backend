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
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.registerUser).Methods("POST")
	router.HandleFunc("/auth/verify-account", h.verifyAccount).Methods("POST")
	router.HandleFunc("/auth/login", h.login).Methods("POST")
}

func (h *AuthHandler) registerUser(w http.ResponseWriter, r *http.Request) {
	var dto RegisterUserRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Per-request instantiation
	unitOfWork := uow.NewUnitOfWork(h.db)
	userRepository := persistance.NewUserRepository(unitOfWork)
	otpRepository := persistance.NewOtpRepository(unitOfWork)
	loginAttemptRepository := persistance.NewLoginAttemptRepository(unitOfWork)
	hashingService := services.NewHashingService()
	notificationService := services.NewLogNotificationService()
	authOtpService := services.NewAuthOtpService(otpRepository)
	userService := services.NewUserService(userRepository, loginAttemptRepository, hashingService)
	registerUserUseCase := usecases.NewRegisterUser(unitOfWork, userService, authOtpService, hashingService, notificationService)

	input := usecases.RegisterUserInput{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
	}

	_, err := registerUserUseCase.Execute(r.Context(), input)
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

	// Per-request instantiation
	unitOfWork := uow.NewUnitOfWork(h.db)
	userRepository := persistance.NewUserRepository(unitOfWork)
	otpRepository := persistance.NewOtpRepository(unitOfWork)
	loginAttemptRepository := persistance.NewLoginAttemptRepository(unitOfWork)
	hashingService := services.NewHashingService()
	authOtpService := services.NewAuthOtpService(otpRepository)
	userService := services.NewUserService(userRepository, loginAttemptRepository, hashingService)
	verifyAccountUseCase := usecases.NewVerifyAccount(unitOfWork, userService, authOtpService, userRepository)

	input := usecases.VerifyAccountInput{
		Email:   dto.Email,
		OTPCode: dto.OTPCode,
	}

	err := verifyAccountUseCase.Execute(r.Context(), input)
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

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var dto LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Per-request instantiation
	unitOfWork := uow.NewUnitOfWork(h.db)
	userRepository := persistance.NewUserRepository(unitOfWork)
	loginAttemptRepository := persistance.NewLoginAttemptRepository(unitOfWork)
	hashingService := services.NewHashingService()
	// IMPORTANT: You should get these from your settings/config in a real app
	jwtService := services.NewJWTService("your-super-secret-key", 30)
	userService := services.NewUserService(userRepository, loginAttemptRepository, hashingService)
	loginUserUseCase := usecases.NewLoginUser(unitOfWork, userService, jwtService)

	input := usecases.LoginUserInput{
		Email:    dto.Email,
		Password: dto.Password,
	}

	result, err := loginUserUseCase.Execute(r.Context(), input)
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
