package services

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/repository"
	"encoding/base64"
	"log/slog"
	"time"
)

type AdminService struct {
	oauth2            OAuth2Interface
	adminRepository   repository.AdminRepository
	fileLogRepository repository.FileLogsRepository
	cryptoUtil        CryptographicInterface
}

func NewAdminService(oauth2 OAuth2Interface, adminRepo repository.AdminRepository, fileLogsRepository repository.FileLogsRepository, cryptoUtil CryptographicInterface) AdminInterface {
	return &AdminService{
		oauth2:            oauth2,
		adminRepository:   adminRepo,
		fileLogRepository: fileLogsRepository,
		cryptoUtil:        cryptoUtil,
	}
}

func (a *AdminService) Login(ctx context.Context, username string, password string) (string, error) {
	if (username == "") || (password == "") {
		return "", model.ErrInvalidInput
	}

	admin, err := a.adminRepository.GetByUsername(ctx, username)
	if err != nil {
		return "", model.AdminErrWrongCredentials
	}

	// Add constant-time delay to prevent timing attacks
	defer func() {
		time.Sleep(100 * time.Millisecond) // Constant delay
	}()

	result, err := a.oauth2.TokenRequest(ctx, &model.TokenRequest{
		ClientId:     admin.ClientID,
		ClientSecret: password,
		GrantType:    "client_credentials",
		Scope:        "offline",
	})
	if err != nil {
		return "", model.AdminErrWrongCredentials
	}

	return result.AccessToken, nil
}

func (a *AdminService) Register(ctx context.Context, username string, password string) (string, error) {
	if (username == "") || (password == "") {
		return "", model.ErrInvalidInput
	}

	salt, err := a.cryptoUtil.GenerateKey()
	if err != nil {
		slog.Error("Salt generation failed", slog.Any("error", err))
		return "", err
	}

	admin, err := a.oauth2.CreateClient(ctx, &model.ApplicationRequest{
		ClientName:   username,
		ClientSecret: password,
		GrantTypes:   []string{"client_credentials", "refresh_token"},
		Scopes:       []string{"offline"},
		ClientUri:    username,
		RedirectUris: []string{"http://localhost:3000"},
	})

	if err != nil {
		return "", err
	}

	input := username + ":" + admin.ClientId
	secret, err := a.encryptSecret(input, salt, admin.ClientSecret)
	if err != nil {
		slog.Error("Hashing failed", slog.Any("error", err))
		err = a.oauth2.DeleteClient(ctx, admin.ClientId)
		if err != nil {
			slog.Error("Failed to delete OAuth2 client after hashing failure", slog.Any("error", err))
			// ADD AUDIT LOG HERE
			a.fileLogRepository.Create(context.Background(), &entity.FileLogs{
				FileID:    "",
				ActorType: "system",
				ActorID:   admin.ClientId,
				Action:    "delete",
				IP:        helper.GetClientIP(ctx),
				UserAgent: helper.GetUserAgent(ctx),
				Metadata:  map[string]interface{}{"reason": "failed to delete OAuth2 client after hashing failure"},
			})
		}
		return "", err
	}

	err = a.adminRepository.Create(ctx, &entity.Admins{
		ID:       helper.GenerateCustomUUID().String(),
		Username: username,
		ClientID: admin.ClientId,
		Secret:   secret,
		Salt:     salt,
	})

	if err != nil {
		slog.Error("Failed to create admin", slog.Any("error", err))
		err = a.oauth2.DeleteClient(ctx, admin.ClientId)
		if err != nil {
			slog.Error("Failed to delete OAuth2 client after admin creation failure", slog.Any("error", err))
			// ADD AUDIT LOG HERE
			a.fileLogRepository.Create(context.Background(), &entity.FileLogs{
				FileID:    "",
				ActorType: "system",
				ActorID:   admin.ClientId,
				Action:    "delete",
				IP:        helper.GetClientIP(ctx),
				UserAgent: helper.GetUserAgent(ctx),
				Metadata:  map[string]interface{}{"reason": "failed to delete OAuth2 client after admin creation failure"},
			})
		}
		return "", err
	}
	return "Success creating admin with username: " + admin.ClientName, nil
}

func (a *AdminService) RefreshToken(ctx context.Context, adminID, accessToken string) (string, error) {
	admin, err := a.adminRepository.GetByClientID(ctx, adminID)
	if err != nil {
		return "", err
	}

	salt := admin.Salt
	input := admin.Username + ":" + admin.ClientID
	adminSecret, err := a.decryptSecret(input, salt, admin.Secret)
	if err != nil {
		return "", err
	}

	defer secureKeyString(adminSecret)() // securely erase adminSecret from memory

	_, err = a.oauth2.RevokeToken(ctx, admin.ClientID, adminSecret, accessToken)
	if err != nil {
		return "", err
	}

	result, err := a.oauth2.TokenRequest(ctx, &model.TokenRequest{
		ClientId:     admin.ClientID,
		ClientSecret: adminSecret,
		GrantType:    "client_credentials",
		Scope:        "offline",
	})
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func (a *AdminService) RevokeToken(ctx context.Context, adminID, accessToken string) (string, error) {
	admin, err := a.adminRepository.GetByClientID(ctx, adminID)
	if err != nil {
		return "", err
	}

	salt := admin.Salt
	input := admin.Username + ":" + admin.ClientID
	adminSecret, err := a.decryptSecret(input, salt, admin.Secret)
	if err != nil {
		return "", err
	}

	defer secureKeyString(adminSecret)() // securely erase adminSecret from memory

	_, err = a.oauth2.RevokeToken(ctx, admin.ClientID, adminSecret, accessToken)
	if err != nil {
		return "", err
	}
	return "Success revoking token", nil
}

func (a *AdminService) UpdateUsername(ctx context.Context, adminID, newUsername string) error {
	admin, err := a.adminRepository.GetByClientID(ctx, adminID)
	if err != nil {
		return err
	}

	_, err = a.oauth2.UpdateClient(ctx, admin.ClientID, "replace", "client_name", newUsername)
	if err != nil {
		return err
	}

	admin.Username = newUsername
	return a.adminRepository.Update(ctx, admin)
}

func (a *AdminService) UpdatePassword(ctx context.Context, adminID, newPassword string) error {
	admin, err := a.adminRepository.GetByClientID(ctx, adminID)
	if err != nil {
		return err
	}

	salt := admin.Salt
	input := admin.Username + ":" + admin.ClientID
	adminSecret, err := a.encryptSecret(input, salt, newPassword)
	if err != nil {
		return err
	}

	_, err = a.oauth2.UpdateClient(ctx, admin.ClientID, "replace", "client_secret", newPassword)
	if err != nil {
		return err
	}
	admin.Secret = adminSecret
	defer secureKeyString(adminSecret)() // securely erase newPassword and adminSecret from memory
	return a.adminRepository.Update(ctx, admin)
}

func (a *AdminService) DeleteAdmin(ctx context.Context, request string) (string, error) {

	if request == "" {
		return "", model.ErrInvalidInput
	}
	admin, err := a.adminRepository.GetByID(ctx, request)
	if err != nil {
		return "", model.AdminErrNotFound
	}
	err = a.adminRepository.Delete(ctx, admin.ID)
	if err != nil {
		return "", err
	}
	err = a.oauth2.DeleteClient(ctx, admin.ClientID)
	if err != nil {
		return "", err
	}

	return "Success deleting admin with username: " + admin.Username, nil
}

func (a *AdminService) GetAdminList(ctx context.Context, offset int, limit int, sortBy, order string) (*[]model.AdminResponse, error) {
	// Validate sort parameters to prevent SQL injection
	sortBy, order = helper.ValidateSortParams(sortBy, order, helper.AllowedAdminSortFields)

	admins, err := a.adminRepository.GetList(ctx, offset, limit, sortBy, order)
	if err != nil {
		return nil, err
	}

	var adminResponses []model.AdminResponse
	for _, admin := range *admins {
		adminResponses = append(adminResponses, model.AdminResponse{
			ID:        admin.ID,
			Username:  admin.Username,
			CreatedAt: admin.CreatedAt.Format(time.RFC3339),
			UpdatedAt: admin.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &adminResponses, nil
}

func (a *AdminService) encryptSecret(input, salt, secret string) (string, error) {

	key, err := a.cryptoUtil.KeyDerivationFunction(input, []byte(salt))
	if err != nil {
		return "", err
	}
	keyBase64 := base64.StdEncoding.EncodeToString(key)
	defer secureKeyString(keyBase64)() // securely erase key from memory
	return a.cryptoUtil.EncryptString(keyBase64, secret)
}

func (a *AdminService) decryptSecret(input, salt, encryptedSecret string) (string, error) {
	key, err := a.cryptoUtil.KeyDerivationFunction(input, []byte(salt))
	if err != nil {
		return "", err
	}
	keyBase64 := base64.StdEncoding.EncodeToString(key)
	defer secureKeyString(keyBase64)() // securely erase key from memory
	return a.cryptoUtil.DecryptString(keyBase64, encryptedSecret)
}
