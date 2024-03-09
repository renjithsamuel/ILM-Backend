package model

import (
	"time"
)

type RoleType string

const (
	Patrons   RoleType = "patrons"
	Librarian RoleType = "librarian"
)

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

type RegisterUserRequest struct {
	Email           string   `json:"email" binding:"required,email"`
	ProfileImageUrl string   `json:"profileImageUrl" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Role            RoleType `json:"role" binding:"required,oneof=librarian patrons"`
}

type LoginUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}
