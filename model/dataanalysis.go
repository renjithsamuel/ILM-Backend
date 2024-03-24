package model

// GetBooksByApproximateDemandRequest
type GetBooksByApproximateDemandRequest struct {
	Page  uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit uint32 `json:"limit" form:"limit" binding:"required,min=5"`
}

// GetRecommendedBooksForUserRequest
type GetRecommendedBooksForUserRequest struct {
	Page  uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit uint32 `json:"limit" form:"limit" binding:"required,min=5"`
}

// Define point weights for different factors
const (
	RatingWeight   = 5
	ReviewWeight   = 1
	ViewWeight     = 1
	WishlistWeight = 1
)
