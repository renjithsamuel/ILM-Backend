package model

// GetAllBooksRequest
type SearchRequest struct {
	Page       *uint32            `json:"page" form:"page" binding:"required,min=1"`
	Limit      *uint32            `json:"limit" form:"limit" binding:"required,min=5"`
	SortBy     *string            `json:"sortBy" form:"sortBy" binding:"required"`
	OrderBy    *string            `json:"orderBy" form:"orderBy" binding:"required"`
	SearchText *string            `json:"searchText" form:"searchText" binding:"required"`
	SearchBy   *string            `json:"searchBy" form:"searchBy" binding:"required"`
	Type       *SearchRequestType `json:"type" form:"type" binding:"required,oneof=user book checkout review"`
}

type SearchRequestType string

var (
	// SearchRequestTypeUser
	SearchRequestTypeUser SearchRequestType = "user"
	// SearchRequestTypeBook
	SearchRequestTypeBook SearchRequestType = "book"
	// SearchRequestTypeCheckout
	SearchRequestTypeCheckout SearchRequestType = "checkout"
	// SearchRequestTypeReview
	SearchRequestTypeReview SearchRequestType = "review"
)
