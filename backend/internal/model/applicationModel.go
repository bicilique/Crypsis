package model

type ApplicationRequest struct {
	ClientId     string   `json:"clientId" gorm:"size:255;unique;not null" validate:"required"`
	ClientSecret string   `json:"clientSecret" gorm:"size:255;unique;not null" validate:"required"`
	ClientName   string   `json:"clientName" gorm:"size:255;unique;not null" validate:"required"`
	ClientUri    string   `json:"clientUri" gorm:"size:255;unique;not null" validate:"required"`
	GrantTypes   []string `json:"grantTypes" gorm:"size:255;unique;not null" validate:"required"`
	RedirectUris []string `json:"redirectUris" gorm:"size:255;unique;not null" validate:"required"`
	Scopes       []string `json:"scopes" gorm:"size:255;unique;not null" validate:"required"`
}
