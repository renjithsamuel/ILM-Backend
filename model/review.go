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

// RatingEntity contains number of ratings against each ratings point
type RatingEntity struct {
	OneStar   uint32 `json:"1star"`
	TwoStar   uint32 `json:"2star"`
	ThreeStar uint32 `json:"3star"`
	FourStar  uint32 `json:"4star"`
	FiveStar  uint32 `json:"5star"`
}

// GetAverageRatingResponse ..
type GetAverageRatingResponse struct {
	Rating       *float64     `json:"rating"`
	TotalReviews uint32       `json:"totalReviews"`
	RatingEntity RatingEntity `json:"ratingEntity"`
}

// ReviewSort
type ReviewSort struct {
	Page   *uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit  *uint32 `json:"limit" form:"limit" binding:"required,min=5"`
	SortBy *string `json:"sortBy" form:"sortBy" binding:"required"`
}
