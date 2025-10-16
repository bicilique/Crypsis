package model

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateAdminUsernameRequest struct {
	Username string `json:"username" validate:"required"`
}

type UpdateAdminPasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

type RekeyRequest struct {
	KeyUID string `json:"keyUID" validate:"required"`
}

type AdminResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
