package usecases

import (
	"context"
	"fmt"

	"github.com/fraineri/plexo_backend/internal/app_info/entities"
	"github.com/fraineri/plexo_backend/internal/app_info/services"
	"github.com/fraineri/plexo_backend/internal/core/persistance/uow"
)

type DisableAppUsecase interface {
	Execute(ctx context.Context) (entities.AppInfo, error)
}

type DisableApp struct {
	uow            uow.UnitOfWork
	appInfoService services.AppInfoService
}

func NewDisableApp(
	uow uow.UnitOfWork,
	appInfoService services.AppInfoService) DisableAppUsecase {
	return &DisableApp{
		uow:            uow,
		appInfoService: appInfoService,
	}
}

func (g *DisableApp) Execute(ctx context.Context) (entities.AppInfo, error) {
	var result entities.AppInfo

	err := g.uow.Do(ctx, func() error {
		appInfo, err := g.appInfoService.GetAppInfo(ctx)

		if err != nil {
			return err
		}

		if appInfo.Name == "" {
			return fmt.Errorf("app info not found")
		}

		appInfo.Status = "DISABLED"
		err = g.appInfoService.DisableApp(ctx, appInfo)
		if err != nil {
			return fmt.Errorf("failed to update app info: %w", err)
		}
		result = appInfo
		return nil
	})

	if err != nil {
		return entities.AppInfo{}, fmt.Errorf("failed to disable app: %w", err)
	}
	return result, nil
}
