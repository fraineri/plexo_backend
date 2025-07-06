package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/app_info/entities"
	"github.com/fraineri/plexo_backend/internal/app_info/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type EnableAppUsecase interface {
	Execute(ctx context.Context) (entities.AppInfo, error)
}

type EnableApp struct {
	uow            uow.UnitOfWork
	appInfoService services.AppInfoService
}

func NewEnableApp(
	uow uow.UnitOfWork,
	appInfoService services.AppInfoService) EnableAppUsecase {
	return &EnableApp{
		uow:            uow,
		appInfoService: appInfoService,
	}
}

func (g *EnableApp) Execute(ctx context.Context) (entities.AppInfo, error) {
	var result entities.AppInfo

	err := g.uow.Do(ctx, func() error {
		appInfo, err := g.appInfoService.GetAppInfo(ctx)

		if err != nil {
			return err
		}

		if appInfo.Name == "" {
			return fmt.Errorf("app info not found")
		}

		appInfo.Status = "ENABLED"
		err = g.appInfoService.EnableApp(ctx, appInfo)
		if err != nil {
			return fmt.Errorf("failed to update app info: %w", err)
		}
		result = appInfo
		return nil
	})

	if err != nil {
		return entities.AppInfo{}, fmt.Errorf("failed to enable app: %w", err)
	}
	return result, nil
}
