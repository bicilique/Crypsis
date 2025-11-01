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

// // Mock OAuth2Interface
// type MockOAuth2 struct {
// 	mock.Mock
// }

// func (m *MockOAuth2) CreateClient(ctx context.Context, req *model.ApplicationRequest) (*model.OAuth2ClientResponse, error) {
// 	args := m.Called(ctx, req)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*model.OAuth2ClientResponse), args.Error(1)
// }

// func (m *MockOAuth2) GetClient(ctx context.Context, clientId string) (*model.OAuth2ClientResponse, error) {
// 	args := m.Called(ctx, clientId)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*model.OAuth2ClientResponse), args.Error(1)
// }

// func (m *MockOAuth2) UpdateClient(ctx context.Context, clientId, op, path string, value interface{}) (*model.OAuth2ClientResponse, error) {
// 	args := m.Called(ctx, clientId, op, path, value)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*model.OAuth2ClientResponse), args.Error(1)
// }

// func (m *MockOAuth2) DeleteClient(ctx context.Context, clientId string) error {
// 	args := m.Called(ctx, clientId)
// 	return args.Error(0)
// }

// func (m *MockOAuth2) TokenRequest(ctx context.Context, req *model.TokenRequest) (*model.TokenResponse, error) {
// 	args := m.Called(ctx, req)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*model.TokenResponse), args.Error(1)
// }

// func (m *MockOAuth2) RevokeToken(ctx context.Context, clientId, clientSecret, token string) (string, error) {
// 	args := m.Called(ctx, clientId, clientSecret, token)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockOAuth2) IntrospectToken(ctx context.Context, token string) (*model.IntrospectResponse, error) {
// 	args := m.Called(ctx, token)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*model.IntrospectResponse), args.Error(1)
// }

// // Mock AdminRepository
// type MockAdminRepository struct {
// 	mock.Mock
// }

// func (m *MockAdminRepository) Create(ctx context.Context, admin *entity.Admins) error {
// 	args := m.Called(ctx, admin)
// 	return args.Error(0)
// }

// func (m *MockAdminRepository) GetByUsername(ctx context.Context, username string) (*entity.Admins, error) {
// 	args := m.Called(ctx, username)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Admins), args.Error(1)
// }

// func (m *MockAdminRepository) GetByClientID(ctx context.Context, clientID string) (*entity.Admins, error) {
// 	args := m.Called(ctx, clientID)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Admins), args.Error(1)
// }

// func (m *MockAdminRepository) GetByID(ctx context.Context, id string) (*entity.Admins, error) {
// 	args := m.Called(ctx, id)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*entity.Admins), args.Error(1)
// }

// func (m *MockAdminRepository) Update(ctx context.Context, admin *entity.Admins) error {
// 	args := m.Called(ctx, admin)
// 	return args.Error(0)
// }

// func (m *MockAdminRepository) Delete(ctx context.Context, id string) error {
// 	args := m.Called(ctx, id)
// 	return args.Error(0)
// }

// func (m *MockAdminRepository) GetList(ctx context.Context, offset, limit int, sortBy, order string) (*[]entity.Admins, error) {
// 	args := m.Called(ctx, offset, limit, sortBy, order)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*[]entity.Admins), args.Error(1)
// }

// func (m *MockAdminRepository) IsAdmin(ctx context.Context, clientID string) bool {
// 	args := m.Called(ctx, clientID)
// 	return args.Bool(0)
// }

// // Mock CryptographicInterface
// type MockCryptographic struct {
// 	mock.Mock
// }

// func (m *MockCryptographic) GenerateKey() (string, error) {
// 	args := m.Called()
// 	return args.String(0), args.Error(1)
// }

// func (m *MockCryptographic) EncryptString(key, plaintext string) (string, error) {
// 	args := m.Called(key, plaintext)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockCryptographic) DecryptString(key, ciphertext string) (string, error) {
// 	args := m.Called(key, ciphertext)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockCryptographic) EncryptFile(key string, plaintext []byte) ([]byte, error) {
// 	args := m.Called(key, plaintext)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]byte), args.Error(1)
// }

// func (m *MockCryptographic) DecryptFile(key string, ciphertext []byte) ([]byte, error) {
// 	args := m.Called(key, ciphertext)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]byte), args.Error(1)
// }

// func (m *MockCryptographic) HashFile(method string, data []byte) (string, error) {
// 	args := m.Called(method, data)
// 	return args.String(0), args.Error(1)
// }

// func (m *MockCryptographic) CompareHashFile(method string, data []byte, hash string) bool {
// 	args := m.Called(method, data, hash)
// 	return args.Bool(0)
// }

// func (m *MockCryptographic) KeyDerivationFunction(password string, salt []byte) ([]byte, error) {
// 	args := m.Called(password, salt)
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).([]byte), args.Error(1)
// }

// func (m *MockCryptographic) ImportRawKeyAsBase64(keyBytes []byte) (string, error) {
// 	args := m.Called(keyBytes)
// 	return args.String(0), args.Error(1)
// }

// // Tests
// func TestAdminService_Login_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()
// 	username := "admin"
// 	password := "password123"

// 	admin := &entity.Admins{
// 		ID:       "admin-id",
// 		Username: username,
// 		ClientID: "client-id",
// 	}

// 	tokenResponse := &model.TokenResponse{
// 		AccessToken: "access-token-123",
// 	}

// 	mockAdminRepo.On("GetByUsername", ctx, username).Return(admin, nil)
// 	mockOAuth2.On("TokenRequest", ctx, mock.AnythingOfType("*model.TokenRequest")).Return(tokenResponse, nil)

// 	token, err := adminService.Login(ctx, username, password)

// 	assert.NoError(t, err)
// 	assert.Equal(t, "access-token-123", token)
// 	mockAdminRepo.AssertExpectations(t)
// 	mockOAuth2.AssertExpectations(t)
// }

// func TestAdminService_Login_EmptyCredentials(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()

// 	_, err := adminService.Login(ctx, "", "password")
// 	assert.Equal(t, model.ErrInvalidInput, err)

// 	_, err = adminService.Login(ctx, "username", "")
// 	assert.Equal(t, model.ErrInvalidInput, err)
// }

// func TestAdminService_Login_AdminNotFound(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()
// 	username := "nonexistent"
// 	password := "password123"

// 	mockAdminRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("admin not found"))

// 	_, err := adminService.Login(ctx, username, password)

// 	assert.Equal(t, model.AdminErrWrongCredentials, err)
// 	mockAdminRepo.AssertExpectations(t)
// }

// func TestAdminService_Register_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()
// 	username := "newadmin"
// 	password := "password123"

// 	oauthResponse := &model.OAuth2ClientResponse{
// 		ClientId:     "client-id-123",
// 		ClientSecret: password,
// 		ClientName:   username,
// 	}

// 	mockCrypto.On("GenerateKey").Return("salt-key-base64", nil)
// 	mockOAuth2.On("CreateClient", ctx, mock.AnythingOfType("*model.ApplicationRequest")).Return(oauthResponse, nil)
// 	mockCrypto.On("KeyDerivationFunction", mock.Anything, mock.Anything).Return([]byte("derived-key"), nil)
// 	mockCrypto.On("EncryptString", mock.Anything, password).Return("encrypted-secret", nil)
// 	mockAdminRepo.On("Create", ctx, mock.AnythingOfType("*entity.Admins")).Return(nil)

// 	result, err := adminService.Register(ctx, username, password)

// 	assert.NoError(t, err)
// 	assert.Contains(t, result, "Succes creating admin")
// 	mockCrypto.AssertExpectations(t)
// 	mockOAuth2.AssertExpectations(t)
// 	mockAdminRepo.AssertExpectations(t)
// }

// func TestAdminService_Register_EmptyCredentials(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()

// 	_, err := adminService.Register(ctx, "", "password")
// 	assert.Equal(t, model.ErrInvalidInput, err)

// 	_, err = adminService.Register(ctx, "username", "")
// 	assert.Equal(t, model.ErrInvalidInput, err)
// }

// func TestAdminService_DeleteAdmin_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()
// 	adminID := "admin-id-123"

// 	admin := &entity.Admins{
// 		ID:       adminID,
// 		Username: "testadmin",
// 		ClientID: "client-id-123",
// 	}

// 	mockAdminRepo.On("GetByID", ctx, adminID).Return(admin, nil)
// 	mockAdminRepo.On("Delete", ctx, adminID).Return(nil)
// 	mockOAuth2.On("DeleteClient", ctx, admin.ClientID).Return(nil)

// 	result, err := adminService.DeleteAdmin(ctx, adminID)

// 	assert.NoError(t, err)
// 	assert.Contains(t, result, "Succes deleting admin")
// 	mockAdminRepo.AssertExpectations(t)
// 	mockOAuth2.AssertExpectations(t)
// }

// func TestAdminService_GetAdminList_Success(t *testing.T) {
// 	mockOAuth2 := new(MockOAuth2)
// 	mockAdminRepo := new(MockAdminRepository)
// 	mockCrypto := new(MockCryptographic)

// 	adminService := services.NewAdminService(mockOAuth2, mockAdminRepo, mockCrypto)

// 	ctx := context.Background()
// 	offset := 0
// 	limit := 10

// 	admins := &[]entity.Admins{
// 		{
// 			ID:        "admin-1",
// 			Username:  "admin1",
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 		{
// 			ID:        "admin-2",
// 			Username:  "admin2",
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}

// 	mockAdminRepo.On("GetList", ctx, offset, limit, "created_at", "desc").Return(admins, nil)

// 	result, err := adminService.GetAdminList(ctx, offset, limit)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Len(t, *result, 2)
// 	mockAdminRepo.AssertExpectations(t)
// }
