package model

type SignUpRequest struct {
	Key         string `json:"-"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Username    string `json:"username" binding:"required"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	DisplayName string `json:"display_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
}

type GetUsersRequest struct {
	IDs      string `form:"ids"`
	ID       string `form:"id"`
	Username string `form:"username"`
	Email    string `form:"email"`
	Status   string `form:"status"`
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
