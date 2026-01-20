package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email,omitempty"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Theme        string    `json:"theme,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type TokenRequest struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeThemeRequest struct {
	Theme string `json:"theme"`
}
