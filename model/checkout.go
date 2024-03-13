package model

import "time"

// CheckoutTicket represents the checkout ticket entity
type CheckoutTicket struct {
	ID           string    `json:"ID"`
	BookID       string    `json:"bookID" binding:"required"`
	UserID       string    `json:"userID" binding:"required"`
	IsCheckedOut bool      `json:"isCheckedOut" binding:"required"`
	IsReturned   bool      `json:"isReturned" binding:"required"`
	NumberOfDays int64     `json:"numberOfDays"`
	FineAmount   float64   `json:"fineAmount"`
	ReservedOn   time.Time `json:"reservedOn"`
	CheckedOutOn time.Time `json:"checkedOutOn"`
	ReturnedDate time.Time `json:"returnedDate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// CheckoutTicketResponse 
type CheckoutTicketResponse struct {
	ID           string    `json:"ID"`
	BookID       string    `json:"bookID" binding:"required"`
	UserID       string    `json:"userID" binding:"required"`
	IsCheckedOut bool      `json:"isCheckedOut" binding:"required"`
	IsReturned   bool      `json:"isReturned" binding:"required"`
	NumberOfDays int64     `json:"numberOfDays"`
	FineAmount   float64   `json:"fineAmount"`
	ReservedOn   time.Time `json:"reservedOn"`
	CheckedOutOn time.Time `json:"checkedOutOn"`
	ReturnedDate time.Time `json:"returnedDate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Book `json:"book"`
	User `json:"user"`
}

type CreateCheckoutRequest struct {
	BookID       string    `json:"bookID" binding:"required"`
	UserID       string    `json:"userID" binding:"required"`
	NumberOfDays int64     `json:"numberOfDays"`
}

type UpdateCheckoutTicketRequest struct {
	ID           string    `json:"ID" binding:"required"`
	BookID       string    `json:"bookID" binding:"required"`
	UserID       string    `json:"userID" binding:"required"`
	IsCheckedOut bool      `json:"isCheckedOut" binding:"omitempty"`
	IsReturned   bool      `json:"isReturned" binding:"omitempty"`
	NumberOfDays int64     `json:"numberOfDays"`
	FineAmount   float64   `json:"fineAmount"`
	ReservedOn   time.Time `json:"reservedOn"`
	CheckedOutOn time.Time `json:"checkedOutOn"`
	ReturnedDate time.Time `json:"returnedDate"`
}
