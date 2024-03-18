package model

import (
	"time"
)

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

type CreateBookRequest struct {
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
	Views         int64     `json:"views" binding:"required"`
	BooksLeft     int64     `json:"booksLeft" binding:"required"`
	// list
	WishList    []string `json:"wishList"`
	ReviewsList []string `json:"reviewsList"`
	ViewsList   []string `json:"viewsList"`
	// count
	WishlistCount     int64   `json:"wishlistCount" binding:"required"`
	Rating            float64 `json:"rating" binding:"required"`
	ReviewCount       int64   `json:"reviewCount" binding:"required"`
	ApproximateDemand int64   `json:"approximateDemand" binding:"omitempty"`
}

type GetBookByISBNRequest struct {
	ISBN string `json:"isbn" uri:"isbn" binding:"required"`
}

type BookDetailsFrom string

var (
	BookDetailsFromReserved   BookDetailsFrom = "reserved"
	BookDetailsFromPending    BookDetailsFrom = "pending"
	BookDetailsFromCheckedOut BookDetailsFrom = "checkedOut"
	BookDetailsFromWishLists  BookDetailsFrom = "wishLists"
	BookDetailsFromCompleted  BookDetailsFrom = "completed"
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
	InLibrary     bool      `json:"inLibrary" binding:"required"`
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