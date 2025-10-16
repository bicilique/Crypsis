package middlewere

import (
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenIntrospect struct {
	Active   bool   `json:"active"`
	Scope    string `json:"scope"`
	ClientID string `json:"client_id"`
	Sub      string `json:"sub"`
}

type TokenMiddlewareConfig struct {
	HydraAdminURL string
	RequiredScope string
	AdminRepo     repository.AdminRepository
}

func TokenMiddleware(config TokenMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, ok := validateToken(c, config)
		if !ok {
			return
		}
		c.Set("tokenInfo", tokenInfo)
		c.Next()
	}
}

func AdminTokenMiddleware(config TokenMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo, ok := validateToken(c, config)
		if !ok {
			return
		}

		if !config.AdminRepo.IsAdmin(c.Request.Context(), tokenInfo.ClientID) {
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User is not an admin")
			return
		}

		c.Set("tokenInfo", tokenInfo)
		c.Next()
	}
}

func validateToken(c *gin.Context, config TokenMiddlewareConfig) (*TokenIntrospect, bool) {
	token := c.GetHeader("Authorization")
	if token == "" {
		model.JSONErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Missing or invalid Authorization header")
		return nil, false
	}

	if !strings.HasPrefix(token, "Bearer ") {
		model.JSONErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Missing Bearer prefix in Authorization header")
		return nil, false
	}

	tokenInfo, err := introspectToken(token, config.HydraAdminURL)
	if err != nil || !tokenInfo.Active {
		model.JSONErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", "Token is not active")
		return nil, false
	}

	if !hasScope(tokenInfo.Scope, config.RequiredScope) {
		model.JSONErrorResponse(c, http.StatusForbidden, "Insufficient scope", "Token lacks required scope")
		return nil, false
	}

	return tokenInfo, true
}

func introspectToken(token, hydraAdminURL string) (*TokenIntrospect, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	req, err := http.NewRequest("POST", hydraAdminURL, strings.NewReader("token="+token))
	if err != nil {
		fmt.Printf("Error creating introspection request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error during HTTP request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	var tokenInfo TokenIntrospect
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		fmt.Printf("Error decoding response body: %v\n", err)
		return nil, err
	}

	return &tokenInfo, nil
}

func hasScope(tokenScope, requiredScope string) bool {
	scopes := strings.Fields(tokenScope)
	for _, scope := range scopes {
		if scope == requiredScope {
			return true
		}
	}
	return false
}

func GetUserIDFromToken(c *gin.Context) (string, bool) {
	tokenInfo, exists := c.Get("tokenInfo")
	if !exists {
		model.JSONErrorResponse(c, http.StatusInternalServerError, "Token info not found", "")
		return "", false
	}

	info, ok := tokenInfo.(*TokenIntrospect)
	if !ok {
		model.JSONErrorResponse(c, http.StatusInternalServerError, "Token info type assertion failed", "")
		return "", false
	}
	return info.Sub, true
}

func GetAccessTokenFromHeader(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", false
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	token = strings.TrimSpace(token) // just in case of accidental spaces
	if token == "" {
		return "", false
	}

	return token, true
}
