package model

import (
	"time"
)

// Book
type Book struct {
	ID            string    `json:"ID" binding:"required,uuid"`
	ISBN          string    `json:"ISBN" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Author        string    `json:"author" binding:"required"`
	Genre         string    `json:"genre"`
	PublishedDate time.Time `json:"publishedDate" binding:"required"`
	Description   string    `json:"desc"`
	PreviewLink   string    `json:"previewLink"`
	CoverImage    string    `json:"coverImage" binding:"required"`
	ShelfNumber   int64     `json:"shelfNumber" binding:"required"`
	InLibrary     bool      `json:"inLibrary" binding:"required"`
	BooksLeft     int64     `json:"booksLeft" binding:"required"`
	Rating        float64   `json:"rating" binding:"required"`
	// list
	WishList    []string `json:"wishList"`
	ReviewsList []string `json:"reviewsList"`
	ViewsList   []string `json:"viewsList"`
	// count
	Views             int64      `json:"views" binding:"required"`
	WishlistCount     int64      `json:"wishlistCount" binding:"required"`
	ReviewCount       int64      `json:"reviewCount" binding:"required"`
	ApproximateDemand int64      `json:"approximateDemand" binding:"required"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}

// CreateBookRequest
type CreateBookRequest struct {
	ISBN          string    `json:"ISBN" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Author        string    `json:"author" binding:"required"`
	Genre         string    `json:"genre"`
	PublishedDate time.Time `json:"publishedDate" binding:"required"`
	Description   string    `json:"desc"`
	PreviewLink   string    `json:"previewLink"`
	CoverImage    string    `json:"coverImage" binding:"required"`
	ShelfNumber   *int64    `json:"shelfNumber" binding:"omitempty"`
	InLibrary     *bool     `json:"inLibrary" binding:"omitempty"`
	Views         *int64    `json:"views" binding:"omitempty"`
	BooksLeft     *int64    `json:"booksLeft" binding:"omitempty"`
	// list
	WishList    []string `json:"wishList"`
	ReviewsList []string `json:"reviewsList"`
	ViewsList   []string `json:"viewsList"`
	// count
	WishlistCount     *int64   `json:"wishlistCount" binding:"omitempty"`
	Rating            *float64 `json:"rating" binding:"omitempty"`
	ReviewCount       *int64   `json:"reviewCount" binding:"omitempty"`
	ApproximateDemand *int64   `json:"approximateDemand" binding:"omitempty"`
}

// GetBookByISBNRequest
type GetBookByISBNRequest struct {
	ISBN string `json:"isbn" uri:"isbn" binding:"required"`
}

// BookDetailsFrom
type BookDetailsFrom string

var (
	// BookDetailsFromReserved
	BookDetailsFromReserved BookDetailsFrom = "reserved"
	// BookDetailsFromPending
	BookDetailsFromPending BookDetailsFrom = "pending"
	// BookDetailsFromCheckedOut
	BookDetailsFromCheckedOut BookDetailsFrom = "checkedout"
	// BookDetailsFromWishLists
	BookDetailsFromWishLists BookDetailsFrom = "wishlists"
	// BookDetailsFromCompleted
	BookDetailsFromCompleted BookDetailsFrom = "completed"
)

type GetAllBooksByBookDetailsFromRequest struct {
	UserID          string          `json:"userID" uri:"userid" binding:"required,uuid"`
	BookDetailsFrom BookDetailsFrom `json:"bookDetailsFrom" uri:"bookdetailsfrom"`
}

type UpdateBookRequest struct {
	ISBN          string    `json:"ISBN" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Author        string    `json:"author" binding:"required"`
	Genre         string    `json:"genre"`
	PublishedDate time.Time `json:"publishedDate" binding:"required"`
	Description   string    `json:"desc"`
	PreviewLink   string    `json:"previewLink"`
	CoverImage    string    `json:"coverImage" binding:"required"`
	ShelfNumber   *int64    `json:"shelfNumber" binding:"required"`
	InLibrary     *bool     `json:"inLibrary" binding:"omitempty"`
	BooksLeft     *int64    `json:"booksLeft" binding:"required"`
	Rating        *float64  `json:"rating" binding:"required"`
	// list
	WishList    []string `json:"wishList"`
	ReviewsList []string `json:"reviewsList"`
	ViewsList   []string `json:"viewsList"`
	// count
	Views             *int64 `json:"views" binding:"required"`
	WishlistCount     *int64 `json:"wishlistCount" binding:"required"`
	ReviewCount       *int64 `json:"reviewCount" binding:"required"`
	ApproximateDemand *int64 `json:"approximateDemand" binding:"omitempty"`
}

// GetAllBooksRequest
type GetAllBooksRequest struct {
	Page    uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit   uint32 `json:"limit" form:"limit" binding:"required,min=5"`
	SortBy  string `json:"sortBy" form:"sortBy" binding:"required"`
	OrderBy string `json:"orderBy" form:"orderBy" binding:"required"`
}

// GetAllCheckoutData
type GetAllCheckoutData struct {
	Page    uint32 `json:"page" form:"page" binding:"required,min=1"`
	Limit   uint32 `json:"limit" form:"limit" binding:"required,min=5"`
	SortBy  string `json:"sortBy" form:"sortBy" binding:"required"`
	OrderBy string `json:"orderBy" form:"orderBy" binding:"required"`
}

// similar books request
type SimilarBooksRequest struct {
	ISBN string `json:"isbn" uri:"isbn" binding:"required"`
}
