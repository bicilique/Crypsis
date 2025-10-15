package model

import "time"

type OAuth2ClientResponse struct {
	ClientId     string
	ClientSecret string
	ClientName   string
	ClientUri    string
	GrantTypes   []string
	RedirectUris []string
	Scopes       string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type OAuth2TokenResponse struct {
	Active       bool
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
	Scope        string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type TokenRequest struct {
	ClientId     string `json:"clientId" gorm:"size:255;unique;not null" validate:"required"`
	ClientSecret string `json:"clientSecret" gorm:"size:255;unique;not null" validate:"required"`
	GrantType    string `json:"grantType" gorm:"size:255;unique;not null" validate:"required"`
	Scope        string `json:"scope" gorm:"size:255;unique;not null" validate:"required"`
	Code         string `json:"code" gorm:"size:255;unique;not null" validate:"required"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
