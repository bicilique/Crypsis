package services

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/repository"
	"fmt"
	"log/slog"
	"strings"
)

type ApplicationService struct {
	oauth2        OAuth2Interface
	appRepository repository.ApplicationRepository

	fileLogsRepository repository.FileLogsRepository
}

func NewApplicationService(oauth2 OAuth2Interface, appRepo repository.ApplicationRepository, fileLogsRepository repository.FileLogsRepository) ApplicationInterface {
	return &ApplicationService{
		oauth2:             oauth2,
		appRepository:      appRepo,
		fileLogsRepository: fileLogsRepository,
	}
}

func (a *ApplicationService) AddApp(ctx context.Context, appName, appUri, redirectUri string) (*model.AppDetailResponse, error) {
	if appName == "" || appUri == "" || redirectUri == "" {
		return nil, model.ErrInvalidInput
	}

	app, err := a.appRepository.GetByName(ctx, appName)
	if err != nil && !(strings.Contains(err.Error(), "app not found")) {
		return nil, err // handle real errors
	}
	if app != nil {
		return nil, model.ErrAppAlreadyExists // app already exists
	}

	// Create OAuth2 client
	appCred, err := a.oauth2.CreateClient(ctx, &model.ApplicationRequest{
		ClientName:   appName,
		GrantTypes:   []string{"client_credentials", "authorization_code", "refresh_token"},
		Scopes:       []string{"openid", "offline"},
		ClientUri:    appUri,
		RedirectUris: []string{redirectUri},
	})

	if err != nil {
		return nil, err
	}

	// Create app in DB
	err = a.appRepository.Create(ctx, &entity.Apps{
		ID:           helper.GenerateCustomUUID().String(),
		Name:         appName,
		ClientID:     appCred.ClientId,
		ClientSecret: appCred.ClientSecret,
		Uri:          appUri,
		RedirectUri:  redirectUri,
		IsActive:     true,
	})
	if err != nil {
		slog.Warn("Error creating app in DB , deleting client", slog.String("client_id", appCred.ClientId))
		err = a.oauth2.DeleteClient(ctx, appCred.ClientId)
		if err != nil {
			slog.Error("Failed to delete OAuth2 client after app creation failure", slog.Any("error", err))
			a.fileLogsRepository.Create(context.Background(), &entity.FileLogs{
				FileID:    "",
				ActorType: "system",
				ActorID:   appCred.ClientId,
				Action:    "delete",
				IP:        helper.GetClientIP(ctx),
				UserAgent: helper.GetUserAgent(ctx),
				Metadata:  map[string]interface{}{"reason": "failed to delete OAuth2 client after app creation failure"},
			})
		}
		return nil, err
	}

	return &model.AppDetailResponse{
		ID:           appCred.ClientId,
		AppName:      appName,
		ClientID:     appCred.ClientId,
		ClientSecret: appCred.ClientSecret,
		IsActive:     true,
		Uri:          appUri,
		RedirectUri:  redirectUri,
		CreatedAt:    appCred.CreatedAt.String(),
		UpdatedAt:    appCred.UpdatedAt.String(),
	}, nil

}

func (a *ApplicationService) UpdateApp(ctx context.Context, appUID string, appName, appUri, redirectUri string) (*model.AppDetailResponse, error) {
	app, err := a.checkAppExist(ctx, appUID)
	if err != nil {
		return nil, err
	}
	app.Name = appName
	app.Uri = appUri
	app.RedirectUri = redirectUri
	err = a.appRepository.Update(ctx, app)
	if err != nil {
		return nil, err
	}
	return a.GetInfo(ctx, appUID)
}

func (a *ApplicationService) GetInfo(ctx context.Context, appUID string) (*model.AppDetailResponse, error) {
	app, err := a.checkAppExist(ctx, appUID)
	if err != nil {
		return nil, err
	}
	return &model.AppDetailResponse{
		ID:           app.ID,
		AppName:      app.Name,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		IsActive:     app.IsActive,
		Uri:          app.Uri,
		RedirectUri:  app.RedirectUri,
		CreatedAt:    app.CreatedAt.String(),
		UpdatedAt:    app.UpdatedAt.String(),
	}, nil
}

func (a *ApplicationService) ListApps(ctx context.Context, limit, offset int, sortBy, order string) (int64, *[]model.AppResponse, error) {
	// Validate sort parameters to prevent SQL injection
	sortBy, order = helper.ValidateSortParams(sortBy, order, helper.AllowedAppSortFields)

	count, appList, err := a.appRepository.GetListApps(ctx, offset, limit, sortBy, order)
	if err != nil {
		return 0, nil, err
	}
	var appResponse []model.AppResponse
	for _, app := range *appList {
		appResponse = append(appResponse, model.AppResponse{
			ID:       app.ID,
			AppName:  app.Name,
			IsActive: app.IsActive,
		})
	}
	return count, &appResponse, nil
}

func (a *ApplicationService) DeleteApp(ctx context.Context, appUID string) error {
	app, err := a.checkAppExist(ctx, appUID)
	if err != nil {
		return err
	}

	_, err = a.oauth2.GetClient(ctx, app.ClientID)
	if err != nil {
		return err
	}

	err = a.appRepository.Delete(ctx, appUID)
	if err != nil {
		return err
	}

	_, err = a.oauth2.UpdateClient(ctx, app.ClientID, "replace", "grant_types", []string{}) // delete grant types
	if err != nil {
		return err
	}
	return nil
}

func (a *ApplicationService) RecoverApp(ctx context.Context, appUID string) (*string, error) {
	if appUID == "" {
		return nil, model.ErrInvalidInput
	}

	app, err := a.appRepository.GetByID(ctx, appUID)
	if err != nil {
		return nil, err
	} else if app.IsActive {
		return nil, model.ErrAppAlreadyActive
	} else if app == nil {
		return nil, model.ErrAppNotFound
	}

	_, err = a.oauth2.GetClient(ctx, app.ClientID)
	if err != nil {
		return nil, err
	}

	err = a.appRepository.Restore(ctx, appUID)
	if err != nil {
		return nil, err
	}

	_, err = a.oauth2.UpdateClient(ctx, app.ClientID, "replace", "grant_types", []string{"client_credentials", "authorization_code", "refresh_token"}) // give grant types
	if err != nil {
		return nil, err
	}

	messsage := fmt.Sprintf("Application %s recovered successfully", app.Name)
	return &messsage, nil
}

func (a *ApplicationService) RotateSecret(ctx context.Context, appUID string) (*model.AppDetailResponse, error) {
	app, err := a.checkAppExist(ctx, appUID)
	if err != nil {
		return nil, err
	}

	_, err = a.oauth2.GetClient(ctx, app.ClientID)
	if err != nil {
		return nil, err
	}

	err = a.oauth2.DeleteClient(ctx, app.ClientID)
	if err != nil {
		return nil, err
	}

	newClient, err := a.oauth2.CreateClient(ctx, &model.ApplicationRequest{
		ClientName:   app.Name,
		GrantTypes:   []string{"client_credentials", "authorization_code", "refresh_token"},
		Scopes:       []string{"openid"},
		ClientUri:    app.Uri,
		RedirectUris: []string{app.RedirectUri},
	})
	if err != nil {
		return nil, err
	}

	app.ClientID = newClient.ClientId
	app.ClientSecret = newClient.ClientSecret

	err = a.appRepository.Update(ctx, app)
	if err != nil {
		return nil, err
	}

	return &model.AppDetailResponse{
		ID:           app.ID,
		AppName:      app.Name,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		IsActive:     app.IsActive,
		Uri:          app.Uri,
		RedirectUri:  app.RedirectUri,
	}, nil

}

func (a *ApplicationService) checkAppExist(ctx context.Context, appUID string) (*entity.Apps, error) {
	if appUID == "" {
		return nil, model.ErrInvalidInput
	}

	app, err := a.appRepository.GetByID(ctx, appUID)
	if err != nil {
		return nil, err
	} else if !app.IsActive {
		return nil, model.ErrAppNotActive
	} else if app == nil {
		return nil, model.ErrAppNotFound
	}
	return app, nil
}
