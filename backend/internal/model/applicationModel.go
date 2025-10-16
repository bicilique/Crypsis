package model

type AddAppRequest struct {
	Name        string `json:"name" validate:"required"`
	Uri         string `json:"uri" validate:"required,url"`
	RedirectUri string `json:"redirectUri" validate:"required,url"`
}

type ApplicationRequest struct {
	ClientId     string   `json:"clientId" gorm:"size:255;unique;not null" validate:"required"`
	ClientSecret string   `json:"clientSecret" gorm:"size:255;unique;not null" validate:"required"`
	ClientName   string   `json:"clientName" gorm:"size:255;unique;not null" validate:"required"`
	ClientUri    string   `json:"clientUri" gorm:"size:255;unique;not null" validate:"required"`
	GrantTypes   []string `json:"grantTypes" gorm:"size:255;unique;not null" validate:"required"`
	RedirectUris []string `json:"redirectUris" gorm:"size:255;unique;not null" validate:"required"`
	Scopes       []string `json:"scopes" gorm:"size:255;unique;not null" validate:"required"`
}

type AppResponse struct {
	ID       string `json:"id"`
	AppName  string `json:"app_name"`
	IsActive bool   `json:"is_active"`
}

type AppDetailResponse struct {
	ID           string `json:"id"`
	AppName      string `json:"app_name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	IsActive     bool   `json:"is_active"`
	Uri          string `gorm:"type:text; null"`
	RedirectUri  string `gorm:"type:text; null"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
