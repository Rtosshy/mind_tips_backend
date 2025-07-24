package models

import "time"

type User struct {
	ID         int       `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Name       string    `json:"name" db:"name"`
	GoogleID   string    `json:"google_id" db:"google_id"`
	PictureURL string    `json:"picture_url" db:"picture_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type PublicUser struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url"`
}

type UserUpdateRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}
