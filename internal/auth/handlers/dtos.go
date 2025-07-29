package handlers

type RegisterUserRequestDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type VerifyAccountRequestDTO struct {
	Email   string `json:"email"`
	OTPCode string `json:"otp_code"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
}

type ErrorResponseDTO struct {
	Message string `json:"message"`
}
