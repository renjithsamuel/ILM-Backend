package model

import (
	"time"
)

// BookDetails
type BookDetails struct {
	UserID               string     `json:"userID" binding:"omitempty,uuid"`
	PendingBooksCount    int64      `json:"pendingBooksCount"`
	PendingBooksList     []string   `json:"pendingBooksList"`
	CheckedOutBooksCount int64      `json:"checkedOutBooksCount"`
	CheckedOutBookList   []string   `json:"checkedOutBookList"`
	ReservedBooksCount   int64      `json:"reservedBooksCount"`
	ReservedBookList     []string   `json:"reservedBookList"`
	CompletedBooksCount  int64      `json:"completedBooksCount"`
	CompletedBooksList   []string   `json:"completedBooksList"`
	FavoriteGenres       []string   `json:"favoriteGenres"`
	WishlistBooks        []string   `json:"wishlistBooks"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
}

// reserved - books that are reseved but yet to be checked out
// checkedOut - books that are checked out
// pending - books that are checkedout but are yet to be added review or returned
// completed - books that are both added review and returned
