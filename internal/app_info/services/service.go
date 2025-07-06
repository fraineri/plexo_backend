package services

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/app_info/entities"
	"github.com/fraineri/plexo_backend/internal/app_info/persistance"
)

type AppInfo = entities.AppInfo
type Repository = persistance.Repository

type AppInfoService interface {
	GetAppInfo(ctx context.Context) (AppInfo, error)
	DisableApp(ctx context.Context, appInfo AppInfo) error
	EnableApp(ctx context.Context, appInfo AppInfo) error
}

type Service struct {
	appInfoRepository Repository
}

func NewAppInfoService(appInfoRepository Repository) AppInfoService {
	return &Service{
		appInfoRepository: appInfoRepository,
	}
}

func (s *Service) GetAppInfo(ctx context.Context) (AppInfo, error) {
	appInfo := s.appInfoRepository.GetAppInfo(ctx)

	if appInfo.Name == "" {
		return AppInfo{}, fmt.Errorf("app info not found")
	}
	return appInfo, nil
}

func (s *Service) DisableApp(ctx context.Context, appInfo AppInfo) error {
	if appInfo.Name == "" {
		return fmt.Errorf("app info not found")
	}

	appInfo.Status = "DISABLED"
	err := s.appInfoRepository.Update(ctx, appInfo)
	if err != nil {
		return fmt.Errorf("failed to update app info: %w", err)
	}
	return nil
}

func (s *Service) EnableApp(ctx context.Context, appInfo AppInfo) error {
	if appInfo.Name == "" {
		return fmt.Errorf("app info not found")
	}

	appInfo.Status = "ENABLED"
	err := s.appInfoRepository.Update(ctx, appInfo)
	if err != nil {
		return fmt.Errorf("failed to update app info: %w", err)
	}
	return nil
}
