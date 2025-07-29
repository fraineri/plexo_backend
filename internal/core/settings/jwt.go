package settings

type JWTSettings struct {
	SecretKey       string `env:"JWT_SECRET_KEY,required"`
	SessionTimeMins int    `env:"JWT_SESSION_TIME_MINS,required" envDefault:"30"`
}
