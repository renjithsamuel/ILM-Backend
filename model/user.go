package model

import (
	"time"
)

// RoleType
type RoleType string

const (
	Patrons   RoleType = "patrons"
	Librarian RoleType = "librarian"
)

// User
type User struct {
	UserID          string      `json:"userID" binding:"required,uuid"`
	ProfileImageUrl string      `json:"profileImageUrl"`
	Name            string      `json:"name"`
	Email           string      `json:"email" binding:"required,email"`
	Role            RoleType    `json:"role"`
	DateOfBirth     *time.Time  `json:"dateOfBirth" binding:"omitempty"`
	PhoneNumber     *string     `json:"phoneNumber" binding:"omitempty"`
	Address         *string     `json:"address"`
	JoinedDate      time.Time   `json:"joinedDate"`
	Country         *string     `json:"country"`
	Views           *int64      `json:"views"`
	FineAmount      float64     `json:"fineAmount"`
	Password        string      `json:"password"`
	IsPaymentDone   bool        `json:"isPaymentDone"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       *time.Time  `json:"updatedAt"`
	BookDetails     BookDetails `json:"bookDetails,omitempty" binding:"omitempty"`
}

// RegisterUserRequest
type RegisterUserRequest struct {
	Email           string   `json:"email" binding:"required,email"`
	ProfileImageUrl string   `json:"profileImageUrl" binding:"omitempty"`
	Name            string   `json:"name" binding:"required"`
	Role            RoleType `json:"role" binding:"required,oneof=librarian patrons"`
	Password        string   `json:"password" binding:"required,min=8,max=20,validatepassword"`
}

// LoginUserRequest
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20,validatepassword"`
}

// GetAllUsersRequest
type GetAllUsersRequest struct {
	Page    uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit   uint32 `json:"limit" form:"limit" binding:"required,min=5"`
	SortBy  string `json:"sortBy" form:"sortBy" binding:"required"`
	OrderBy string `json:"orderBy" form:"orderBy" binding:"required"`
}

// SortPagination
type SortPagination struct {
	Page    uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit   uint32 `json:"limit" form:"limit" binding:"required,min=5"`
	SortBy  string `json:"sortBy" form:"sortBy" binding:"required"`
	OrderBy string `json:"orderBy" form:"orderBy" binding:"required"`
}
