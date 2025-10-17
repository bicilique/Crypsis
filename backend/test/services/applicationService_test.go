package services_test

// import (
// 	"context"
// 	"crypsis-backend/internal/entity"
// 	"crypsis-backend/internal/model"
// 	"crypsis-backend/internal/services"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // Mock ApplicationRepository
// type MockApplicationRepository struct {
// 	mock.Mock
// }

// func (m *MockApplicationRepository) Create(ctx context.Context, app *entity.Apps) error {
// 	args := m.Called(ctx, app)
// 	return args.Error(0)
// }

// func (m *MockApplicationRepository) GetByName(ctx context.Context, name string) (*entity.Apps, error) {
// 	args := m.Called(ctx, name)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Apps), args.Error(1)
// }

// func (m *MockApplicationRepository) GetByID(ctx context.Context, id string) (*entity.Apps, error) {
// 	args := m.Called(ctx, id)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Apps), args.Error(1)
// }

// func (m *MockApplicationRepository) GetByClientID(ctx context.Context, clientID string) (*entity.Apps, error) {
// 	args := m.Called(ctx, clientID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Apps), args.Error(1)
// }

// func (m *MockApplicationRepository) GetByClientIDLimited(ctx context.Context, clientID string) (*entity.Apps, error) {
// 	args := m.Called(ctx, clientID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Apps), args.Error(1)
// }

// func (m *MockApplicationRepository) Update(ctx context.Context, app *entity.Apps) error {
// 	args := m.Called(ctx, app)
// 	return args.Error(0)
// }

// func (m *MockApplicationRepository) Delete(ctx context.Context, id string) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockApplicationRepository) GetListApps(ctx context.Context, offset, limit int, sortBy, order string) (int64, *[]entity.Apps, error) {
// 	args := m.Called(ctx, offset, limit, sortBy, order)
// 	return args.Get(0).(int64), args.Get(1).(*[]entity.Apps), args.Error(2)
// }

// func (m *MockApplicationRepository) ActivateOrDeactivate(ctx context.Context, id string, isActive bool) error {
// 	args := m.Called(ctx, id, isActive)
// 	return args.Error(0)
// }

// // Tests
// func TestApplicationService_AddApp_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appName := "TestApp"
// 	appUri := "https://testapp.com"
// 	redirectUri := "https://testapp.com/callback"

// 	oauthResponse := &model.OAuth2ClientResponse{
// 		ClientId:     "client-id-123",
// 		ClientSecret: "secret-123",
// 		ClientName:   appName,
// 		CreatedAt:    &time.Time{},
// 		UpdatedAt:    &time.Time{},
// 	}

// 	mockAppRepo.On("GetByName", ctx, appName).Return(nil, errors.New("app not found"))
// 	mockOAuth2.On("CreateClient", ctx, mock.AnythingOfType("*model.ApplicationRequest")).Return(oauthResponse, nil)
// 	mockAppRepo.On("Create", ctx, mock.AnythingOfType("*entity.Apps")).Return(nil)

// 	result, err := appService.AddApp(ctx, appName, appUri, redirectUri)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, appName, result.AppName)
// 	assert.Equal(t, "client-id-123", result.ClientID)
// 	mockAppRepo.AssertExpectations(t)
// 	mockOAuth2.AssertExpectations(t)
// }

// func TestApplicationService_AddApp_EmptyInput(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()

// 	_, err := appService.AddApp(ctx, "", "uri", "redirect")
// 	assert.Equal(t, model.ErrInvalidInput, err)

// 	_, err = appService.AddApp(ctx, "name", "", "redirect")
// 	assert.Equal(t, model.ErrInvalidInput, err)

// 	_, err = appService.AddApp(ctx, "name", "uri", "")
// 	assert.Equal(t, model.ErrInvalidInput, err)
// }

// func TestApplicationService_AddApp_AlreadyExists(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appName := "ExistingApp"

// 	existingApp := &entity.Apps{
// 		ID:   "app-id",
// 		Name: appName,
// 	}

// 	mockAppRepo.On("GetByName", ctx, appName).Return(existingApp, nil)

// 	_, err := appService.AddApp(ctx, appName, "uri", "redirect")

// 	assert.Equal(t, model.ErrAppAlreadyExists, err)
// 	mockAppRepo.AssertExpectations(t)
// }

// func TestApplicationService_GetInfo_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appUID := "app-id-123"

// 	app := &entity.Apps{
// 		ID:           appUID,
// 		Name:         "TestApp",
// 		ClientID:     "client-id",
// 		ClientSecret: "secret",
// 		IsActive:     true,
// 		Uri:          "https://testapp.com",
// 		RedirectUri:  "https://testapp.com/callback",
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	}

// 	mockAppRepo.On("GetByID", ctx, appUID).Return(app, nil)

// 	result, err := appService.GetInfo(ctx, appUID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, appUID, result.ID)
// 	assert.Equal(t, "TestApp", result.AppName)
// 	mockAppRepo.AssertExpectations(t)
// }

// func TestApplicationService_UpdateApp_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appUID := "app-id-123"
// 	newName := "UpdatedApp"
// 	newUri := "https://updated.com"
// 	newRedirect := "https://updated.com/callback"

// 	app := &entity.Apps{
// 		ID:           appUID,
// 		Name:         "OldApp",
// 		ClientID:     "client-id",
// 		ClientSecret: "secret",
// 		IsActive:     true,
// 		Uri:          "https://old.com",
// 		RedirectUri:  "https://old.com/callback",
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	}

// 	// First call for checkAppExist in UpdateApp
// 	mockAppRepo.On("GetByID", ctx, appUID).Return(app, nil).Once()
// 	mockAppRepo.On("Update", ctx, mock.AnythingOfType("*entity.Apps")).Return(nil)
// 	// Second call for GetInfo
// 	mockAppRepo.On("GetByID", ctx, appUID).Return(&entity.Apps{
// 		ID:           appUID,
// 		Name:         newName,
// 		ClientID:     "client-id",
// 		ClientSecret: "secret",
// 		IsActive:     true,
// 		Uri:          newUri,
// 		RedirectUri:  newRedirect,
// 		CreatedAt:    app.CreatedAt,
// 		UpdatedAt:    time.Now(),
// 	}, nil).Once()

// 	result, err := appService.UpdateApp(ctx, appUID, newName, newUri, newRedirect)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, newName, result.AppName)
// 	mockAppRepo.AssertExpectations(t)
// }

// func TestApplicationService_ListApps_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	limit := 10
// 	offset := 0

// 	apps := &[]entity.Apps{
// 		{ID: "app-1", Name: "App1", IsActive: true},
// 		{ID: "app-2", Name: "App2", IsActive: false},
// 	}

// 	mockAppRepo.On("GetListApps", ctx, offset, limit, "created_at", "desc").Return(int64(2), apps, nil)

// 	count, result, err := appService.ListApps(ctx, limit, offset, "created_at", "desc")

// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(2), count)
// 	assert.Len(t, *result, 2)
// 	mockAppRepo.AssertExpectations(t)
// }

// func TestApplicationService_DeleteApp_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appUID := "app-id-123"

// 	app := &entity.Apps{
// 		ID:       appUID,
// 		ClientID: "client-id",
// 	}

// 	oauthClient := &model.OAuth2ClientResponse{
// 		ClientId: "client-id",
// 	}

// 	mockAppRepo.On("GetByID", ctx, appUID).Return(app, nil)
// 	mockOAuth2.On("GetClient", ctx, app.ClientID).Return(oauthClient, nil)
// 	mockAppRepo.On("Delete", ctx, appUID).Return(nil)
// 	mockOAuth2.On("UpdateClient", ctx, app.ClientID, "replace", "grant_types", mock.Anything).Return(nil, nil)

// 	err := appService.DeleteApp(ctx, appUID)

// 	assert.NoError(t, err)
// 	mockAppRepo.AssertExpectations(t)
// 	mockOAuth2.AssertExpectations(t)
// }

// func TestApplicationService_ActivateApp_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appUID := "app-id-123"

// 	app := &entity.Apps{
// 		ID:       appUID,
// 		IsActive: false,
// 	}

// 	mockAppRepo.On("GetByID", ctx, appUID).Return(app, nil)
// 	mockAppRepo.On("ActivateOrDeactivate", ctx, appUID, true).Return(nil)

// 	err := appService.ActivateApp(ctx, appUID)

// 	assert.NoError(t, err)
// 	mockAppRepo.AssertExpectations(t)
// }

// func TestApplicationService_DeactivateApp_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAppRepo := new(MockApplicationRepository)

// 	appService := services.NewApplicationService(mockOAuth2, mockAppRepo)

// 	ctx := context.Background()
// 	appUID := "app-id-123"

// 	app := &entity.Apps{
// 		ID:       appUID,
// 		IsActive: true,
// 	}

// 	mockAppRepo.On("GetByID", ctx, appUID).Return(app, nil)
// 	mockAppRepo.On("ActivateOrDeactivate", ctx, appUID, false).Return(nil)

// 	err := appService.DeactivateApp(ctx, appUID)

// 	assert.NoError(t, err)
// 	mockAppRepo.AssertExpectations(t)
// }
