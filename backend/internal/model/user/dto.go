package user

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type RegisterPayload struct {
	Name     string `json:"name"     validate:"required,min=2,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (p *RegisterPayload) Validate() error { return validate.Struct(p) }

type LoginPayload struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (p *LoginPayload) Validate() error { return validate.Struct(p) }

type UpdatePayload struct {
	Name            *string `json:"name"             validate:"omitempty,min=2,max=100"`
	Email           *string `json:"email"            validate:"omitempty,email"`
	CurrentPassword *string `json:"current_password" validate:"omitempty"`
	NewPassword     *string `json:"new_password"     validate:"omitempty,min=8"`
}

func (p *UpdatePayload) Validate() error { return validate.Struct(p) }

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type GetUserPayload struct{}

func (p *GetUserPayload) Validate() error { return nil }
