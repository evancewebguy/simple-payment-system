package user

import (
	_ "github.com/go-playground/validator/v10"
)

type ProviderType string

const (
	ProviderEmail ProviderType = "email"
)

// SignUpRequest contains the information required for a user to sign up
type SignUpRequest struct {
	FullName string       `json:"full_name" validate:"omitempty"`
	Email    string       `json:"email" validate:"omitempty"`
	Password string       `json:"password" validate:"omitempty"`
	Token    string       `json:"token" validate:"omitempty"`
	Provider ProviderType `json:"provider" validate:"required,oneof=email"`
}

// SignInRequest contains the information required for a user to sign in
type SignInRequest struct {
	Email    string       `json:"email" validate:"omitempty,email"`
	Password string       `json:"password" validate:"omitempty"`
	Token    string       `json:"token" validate:"omitempty"` // Optional, used if social login
	Provider ProviderType `json:"provider" validate:"required,oneof=email"`
}

type InitiateResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Otp      string `json:"otp" validate:"required"`
}

type UserDto struct {
	ID          uint   `json:"id"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" gorm:"uniqueIndex"`
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`
	IsVerified  bool   `json:"is_verified"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User  UserDto              `json:"user"`
	Token RefreshTokenResponse `json:"token"`
}
