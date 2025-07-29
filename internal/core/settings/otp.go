package settings

type OTPSettings struct {
	RateLimitCount   int `env:"OTP_RATE_LIMIT_COUNT" envDefault:"5"`
	RateLimitMinutes int `env:"OTP_RATE_LIMIT_MINUTES" envDefault:"60"`
}
