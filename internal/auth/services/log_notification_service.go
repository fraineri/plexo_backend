package services

import (
	"context"
	"log"
)

type logNotificationService struct{}

func NewLogNotificationService() NotificationService {
	return &logNotificationService{}
}

func (s *logNotificationService) SendOTP(ctx context.Context, email, otp string) error {
	log.Printf("---- OTP Notification ----")
	log.Printf("Recipient: %s", email)
	log.Printf("Code: %s", otp)
	log.Printf("--------------------------")
	return nil
}
