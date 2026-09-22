package models

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email,omitempty"`
	Name         string `json:"name"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	GroupID      *int   `json:"group_id"`

	Group     Group     `json:"group"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegistrationCode struct {
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	GroupID   *int      `json:"group_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type RegCodeRequest struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

type UserSession struct {
	ID        string    `json:"-"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
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
	GroupID  *int   `json:"group_id"`
}

type AuthResponse struct {
	User      User      `json:"user"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeThemeRequest struct {
	Theme string `json:"theme"`
}
