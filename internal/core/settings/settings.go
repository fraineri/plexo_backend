package settings

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"

	"github.com/fraineri/plexo_backend/internal/app_info/settings"
)

type appInfoSettings = settings.Settings

type Settings struct {
	Server   ServerSettings
	Database DatabaseSettings
	Project  ProjectSettings
	AppInfo  appInfoSettings
	JWT      JWTSettings
	OTP      OTPSettings
}

func LoadConfig() (*Settings, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found. Reading from environment variables.")
	}

	settings := &Settings{}

	if err := env.Parse(settings); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return settings, nil
}
