package model

import "time"

// CheckoutTicket represents the checkout ticket entity
type CheckoutTicket struct {
	ID           string    `json:"ID"`
	BookID       string    `json:"bookID"`
	UserID       string    `json:"userID"`
	IsCheckedOut bool      `json:"isCheckedOut"`
	IsReturned   bool      `json:"isReturned"`
	NumberOfDays int64     `json:"numberOfDays"`
	FineAmount   float64   `json:"fineAmount"`
	ReservedOn   time.Time `json:"reservedOn"`
	CheckedOutOn time.Time `json:"checkedOutOn"`
	ReturnedDate time.Time `json:"returnedDate"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
