package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/app_info/entities"
	"github.com/fraineri/plexo_backend/internal/app_info/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type GetAppInfoUsecase interface {
	Execute(ctx context.Context) (entities.AppInfo, error)
}

type GetAppInfo struct {
	uow            uow.UnitOfWork
	appInfoService services.AppInfoService
}

func NewGetAppInfo(
	uow uow.UnitOfWork,
	appInfoService services.AppInfoService) GetAppInfoUsecase {
	return &GetAppInfo{
		uow:            uow,
		appInfoService: appInfoService,
	}
}

func (g *GetAppInfo) Execute(ctx context.Context) (entities.AppInfo, error) {
	var result entities.AppInfo

	err := g.uow.Do(ctx, func() error {
		appInfo, err := g.appInfoService.GetAppInfo(ctx)
		if err != nil {
			return err
		}

		if appInfo.Name == "" {
			return fmt.Errorf("app info not found")
		}
		result = appInfo
		return nil
	})

	if err != nil {
		return entities.AppInfo{}, fmt.Errorf("failed to get app info: %w", err)
	}
	return result, nil
}
