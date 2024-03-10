package model

import "time"

// Review represents the review model
type Review struct {
	ID             string     `json:"ID"`
	BookID         string     `json:"bookID"`
	CheckoutID     string     `json:"checkoutID"`
	UserID         string     `json:"userID"`
	CommentHeading string     `json:"commentHeading" binding:"required"`
	Comment        string     `json:"comment" binding:"required"`
	Rating         float64    `json:"rating" binding:"required,numeric,min=0,max=5"`
	Likes          int64      `json:"likes"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

// CreateReviewRequest represents the request model for creating a review
type CreateReviewRequest struct {
	BookID         string  `json:"bookID" binding:"required,uuid"`
	CheckoutID     string  `json:"checkoutID" binding:"required,uuid"`
	UserID         string  `json:"userID" binding:"required,uuid"`
	CommentHeading string  `json:"commentHeading" binding:"required"`
	Comment        string  `json:"comment" binding:"required"`
	Rating         float64 `json:"rating" binding:"required,numeric,min=0,max=5"`
}

// UpdateReviewRequest represents the request model for updating a review
type UpdateReviewRequest struct {
	ID             string  `json:"ID" binding:"required,uuid"`
	CommentHeading string  `json:"commentHeading" binding:"required"`
	Comment        string  `json:"comment" binding:"required"`
	Rating         float64 `json:"rating" binding:"required"`
	Likes          int64   `json:"likes" binding:"required"`
}
