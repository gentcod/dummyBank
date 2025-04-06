package api

import (
	"time"
)

// Pagination is used for setting limit and offset for api request to the database
type pagination struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

// GetEntityByIdRequest is used to set binding request for uri using uuid
type getEntityByIdUUIDRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

// GetEntityByIdRequest is used to set binding request for uri using uuid
type getEntityByIdRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// UserAccount contains the details of account created with userId
type UserAccount struct {
	Username  string    `json:"username" binding:"required,alphanum"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProfile contains the details of
type UserProfile struct {
	Username          string    `json:"username" binding:"required,alphanum"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}
