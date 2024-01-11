package api

import (
	"time"
)

//UserAccount contains the details of account created with userId
type UserAccount struct {
	Username        string    `json:"username" binding:"required,alphanum"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//
type UserProfile struct {
	Username        string    `json:"username" binding:"required,alphanum"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}