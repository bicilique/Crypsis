package services

import (
	"bytes"
	"context"
	"crypsis-backend/internal/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"

	ory "github.com/ory/client-go"
)

type HydraService struct {
	apiClientAdmin  *ory.APIClient
	apiClientPublic *ory.APIClient
	publicURL       string
}

func NewHydraService(hydraAdminURL, hydraPublicURL string) OAuth2Interface {
	adminConfig := ory.NewConfiguration()
	adminConfig.Servers = ory.ServerConfigurations{
		{
			URL: hydraAdminURL,
		},
	}
	publicConfig := ory.NewConfiguration()
	publicConfig.Servers = ory.ServerConfigurations{
		{
			URL: hydraPublicURL,
		},
	}

	return &HydraService{
		apiClientAdmin:  ory.NewAPIClient(adminConfig),
		apiClientPublic: ory.NewAPIClient(publicConfig),
		publicURL:       hydraPublicURL,
	}
}

func (h *HydraService) CreateClient(ctx context.Context, input *model.ApplicationRequest) (*model.OAuth2ClientResponse, error) {
	newCLient := createHydraClient(input)
	client, _, err := h.apiClientAdmin.OAuth2API.CreateOAuth2Client(ctx).OAuth2Client(newCLient).Execute()
	if err != nil {
		slog.Error("failed to create client", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	slog.Info("client created", slog.String("client_name", client.GetClientName()))
	return &model.OAuth2ClientResponse{
		ClientId:     client.GetClientId(),
		ClientSecret: client.GetClientSecret(),
		ClientName:   client.GetClientName(),
		ClientUri:    client.GetClientUri(),
		GrantTypes:   client.GetGrantTypes(),
		RedirectUris: client.GetRedirectUris(),
		Scopes:       client.GetScope(),
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}, nil
}

func (h *HydraService) GetClient(ctx context.Context, clientId string) (*model.OAuth2ClientResponse, error) {
	client, _, err := h.apiClientAdmin.OAuth2API.GetOAuth2Client(ctx, clientId).Execute()
	if err != nil {
		slog.Error("failed to get client", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get client: %v", err)
	}
	slog.Info("client found", slog.String("client_name", client.GetClientName()))

	if client.ClientSecret != nil {
		fmt.Println("Client response : " + *client.ClientSecret)
	} else {
		fmt.Println("Client response : <nil>")
	}
	return &model.OAuth2ClientResponse{
		ClientId:     client.GetClientId(),
		ClientSecret: client.GetClientSecret(),
		ClientName:   client.GetClientName(),
		ClientUri:    client.GetClientUri(),
		GrantTypes:   client.GetGrantTypes(),
		RedirectUris: client.GetRedirectUris(),
		Scopes:       client.GetScope(),
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}, nil
}

func (h *HydraService) ListClients(ctx context.Context, pageToken string, perPage int64) ([]*model.OAuth2ClientResponse, error) {
	clients, _, err := h.apiClientAdmin.OAuth2API.ListOAuth2Clients(ctx).PageSize(perPage).PageToken(pageToken).Execute()
	if err != nil {
		slog.Error("failed to list clients", slog.Any("error", err))
		return nil, fmt.Errorf("failed to list clients: %v", err)
	}
	slog.Info("clients listed", slog.Int("count", len(clients)))
	var clientsResponse []*model.OAuth2ClientResponse
	for _, client := range clients {
		clientsResponse = append(clientsResponse, &model.OAuth2ClientResponse{
			ClientId:     client.GetClientId(),
			ClientSecret: client.GetClientSecret(),
			ClientName:   client.GetClientName(),
			ClientUri:    client.GetClientUri(),
			GrantTypes:   client.GetGrantTypes(),
			RedirectUris: client.GetRedirectUris(),
			Scopes:       client.GetScope(),
			CreatedAt:    client.CreatedAt,
			UpdatedAt:    client.UpdatedAt,
		})
	}
	return clientsResponse, nil
}

func (h *HydraService) UpdateClient(ctx context.Context, clientName string, operations string, path string, value interface{}) (*model.OAuth2ClientResponse, error) {
	client, _, err := h.apiClientAdmin.OAuth2API.GetOAuth2Client(ctx, clientName).Execute()
	if err != nil {
		slog.Error("failed to get client", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	client, _, err = h.apiClientAdmin.OAuth2API.PatchOAuth2Client(ctx, client.GetClientId()).JsonPatch([]ory.JsonPatch{{Op: operations, Path: "/" + path, Value: value}}).
		Execute()

	if err != nil {
		slog.Error("failed to update client", slog.Any("error", err))
		return nil, fmt.Errorf("failed to update client: %v", err)
	}
	slog.Info("client updated", slog.String("client_name", client.GetClientName()))
	return &model.OAuth2ClientResponse{
		ClientId:     client.GetClientId(),
		ClientSecret: client.GetClientSecret(),
		ClientName:   client.GetClientName(),
		ClientUri:    client.GetClientUri(),
		GrantTypes:   client.GetGrantTypes(),
		RedirectUris: client.GetRedirectUris(),
		Scopes:       client.GetScope(),
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}, nil

}

func (h *HydraService) DeleteClient(ctx context.Context, clientId string) error {
	_, err := h.apiClientAdmin.OAuth2API.DeleteOAuth2Client(ctx, clientId).Execute()
	if err != nil {
		slog.Error("failed to delete client", slog.Any("error", err))
		return fmt.Errorf("failed to delete client: %v", err)
	}
	slog.Warn("client deleted", slog.String("client_id", clientId))
	return nil
}

func (h *HydraService) IntrospectToken(ctx context.Context, token, scope string) (*model.OAuth2TokenResponse, error) {
	tokenResponse, _, err := h.apiClientAdmin.OAuth2API.IntrospectOAuth2Token(ctx).Token(token).Scope(scope).Execute()
	if err != nil {
		slog.Error("failed to introspect token", slog.Any("error", err))
		return nil, fmt.Errorf("failed to introspect token: %v", err)
	}
	// slog.Info("token introspected", slog.String("token_type", tokenResponse.GetTokenType()))
	return &model.OAuth2TokenResponse{
		Active:    tokenResponse.Active,
		TokenType: tokenResponse.GetTokenType(),
		ExpiresIn: int(tokenResponse.GetExp()),
		Scope:     tokenResponse.GetScope(),
	}, nil
}

func (h *HydraService) TokenRequest(ctx context.Context, input *model.TokenRequest) (*model.TokenResponse, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", input.ClientId)
	form.Add("client_secret", input.ClientSecret)
	form.Add("scope", input.Scope)

	fmt.Println("Token Request to Hydra with client_id: " + input.ClientId)
	fmt.Println("Client Secret: " + input.ClientSecret)
	fmt.Println("Public URL: " + h.publicURL)
	fmt.Println("Form Data: " + input.Scope)
	req, err := http.NewRequest("POST", h.publicURL+"/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		slog.Error("failed to create request", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to send request", slog.Any("error", err))
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read response body", slog.Any("error", err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		var tokenResponse model.TokenResponse
		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			slog.Error("failed to parse success response", slog.String("body", string(body)), slog.Any("error", err))
			return nil, fmt.Errorf("failed to parse success response: %w", err)
		}
		return &tokenResponse, nil
	}

	// Handle error responses
	var errorResponse model.ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		slog.Error("failed to parse error response", slog.String("body", string(body)), slog.Any("error", err))
		fmt.Println("Unexpected Error Response:", string(body))
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}
	slog.Error("token request failed", slog.String("error", errorResponse.Error), slog.String("error_description", errorResponse.ErrorDescription))
	return nil, fmt.Errorf("error: %s, description: %s", errorResponse.Error, errorResponse.ErrorDescription)
}

// func (h *HydraOauth2Client) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
// 	// h.apiClientAdmin.OAuth2API.
// }

func (h *HydraService) RevokeToken(ctx context.Context, clientId, clientSecret, refreshToken string) (string, error) {
	_, err := h.apiClientPublic.OAuth2API.RevokeOAuth2Token(ctx).ClientId(clientId).ClientSecret(clientSecret).Token(refreshToken).Execute()
	if err != nil {
		slog.Error("failed to revoke token", slog.Any("error", err))
		return "", fmt.Errorf("failed to revoke token: %v", err)
	}
	return "token from client with id: " + clientId + " revoked", nil
}

func createHydraClient(input *model.ApplicationRequest) ory.OAuth2Client {
	var grantTypes []string
	if input.GrantTypes != nil || len(input.GrantTypes) != 0 {
		for _, v := range input.GrantTypes {
			grantTypes = append(grantTypes, v)
		}
	} else {
		grantTypes = []string{"authorization_code", "refresh_token", "client_credentials"}
	}

	var redirect_uri []string
	if input.RedirectUris != nil || len(input.RedirectUris) != 0 {
		for _, v := range input.RedirectUris {
			redirect_uri = append(redirect_uri, v)
		}
	} else {
		redirect_uri = []string{"http://localhost:3000"}
	}

	var scopes string
	if input.Scopes != nil || len(input.Scopes) != 0 {
		for _, v := range input.Scopes {
			scopes = scopes + " " + v
		}
	} else {
		scopes = "openid offline"
	}

	response_types := []string{"code", "id_token"}
	tokenAUthMethod := "client_secret_post"
	oAuth2Client := *ory.NewOAuth2Client()
	oAuth2Client.SetClientName(input.ClientName)
	oAuth2Client.SetScope(scopes)
	oAuth2Client.SetGrantTypes(grantTypes[:])
	oAuth2Client.SetRedirectUris(redirect_uri[:])
	oAuth2Client.SetResponseTypes(response_types[:])
	oAuth2Client.SetTokenEndpointAuthMethod(tokenAUthMethod)

	if input.ClientId != "" {
		oAuth2Client.SetClientId(input.ClientId)
	}
	if input.ClientSecret != "" {
		oAuth2Client.SetClientSecret(input.ClientSecret)
	}

	return oAuth2Client
}
