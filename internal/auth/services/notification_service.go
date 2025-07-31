package services

import "context"

type NotificationService interface {
	SendOTP(ctx context.Context, email, otp string) error
}
