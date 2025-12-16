package models

import (
	"github.com/google/uuid"
)

type userResponse struct {
	ID uuid.UUID `json:"id"`
}

type SignInResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
	Error        interface{}  `json:"error"`
	RequiresMFA  bool         `json:"needs_mfa"`
}

type Payload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupResponse struct {
	ID uuid.UUID `json:"id"`
}

type SignupResponse struct {
	AccessToken string             `json:"access_token"`
	User        UserSignupResponse `json:"user"`
	RequiresMFA bool               `json:"needs_mfa"`
}
