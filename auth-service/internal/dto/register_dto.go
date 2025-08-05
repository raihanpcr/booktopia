package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type RegisterResponse struct {
	ID    uint   `json:"id" validate:"required" example:"1"`
	Name  string `json:"name" validate:"required" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john.doe@example.com"`
}

type TemplateRegisterResponse struct {
	StatusCode int              `json:"status_code" validate:"required" example:"201"`
	Message    string           `json:"message" validate:"required" example:"Create user success"`
	Data       RegisterResponse `json:"data"`
}