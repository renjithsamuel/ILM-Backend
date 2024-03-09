package model

import (
	"time"
)

type BookGenreType string

// StringToBookGenre converts a string to BookGenreType
func (b BookGenreType) StringToBookGenre(value string) BookGenreType {
    return BookGenreType(value)
}

const (
	BookGenreType_Mystery           BookGenreType = "mystery"
	BookGenreType_Thriller          BookGenreType = "thriller"
	BookGenreType_ScienceFiction    BookGenreType = "science_fiction"
	BookGenreType_Fantasy           BookGenreType = "fantasy"
	BookGenreType_Romance           BookGenreType = "romance"
	BookGenreType_HistoricalFiction BookGenreType = "historical_fiction"
	BookGenreType_Horror            BookGenreType = "horror"
	BookGenreType_NonFiction        BookGenreType = "non_fiction"
	BookGenreType_Biography         BookGenreType = "biography"
	BookGenreType_Poetry            BookGenreType = "poetry"
	BookGenreType_Comedy            BookGenreType = "comedy"
	BookGenreType_Drama             BookGenreType = "drama"
	BookGenreType_Adventure         BookGenreType = "adventure"
	BookGenreType_Children          BookGenreType = "children"
	BookGenreType_YoungAdult        BookGenreType = "young_adult"
	BookGenreType_Science           BookGenreType = "science"
	BookGenreType_SelfHelp          BookGenreType = "self_help"
	BookGenreType_Philosophy        BookGenreType = "philosophy"
	BookGenreType_Travel            BookGenreType = "travel"
	BookGenreType_Cookbooks         BookGenreType = "cookbooks"
	BookGenreType_GraphicNovel      BookGenreType = "graphic_novel"
	BookGenreType_Classic           BookGenreType = "classic"
	BookGenreType_Dystopian         BookGenreType = "dystopian"
	BookGenreType_HistoricalRomance BookGenreType = "historical_romance"
	BookGenreType_Crime             BookGenreType = "crime"
	BookGenreType_Western           BookGenreType = "western"
	BookGenreType_Humor             BookGenreType = "humor"
	BookGenreType_Other             BookGenreType = "other"
)

type BookDetails struct {
	UserID               string          `json:"userID" binding:"omitempty,uuid"`
	PendingBooksCount    int64           `json:"pendingBooksCount"`
	PendingBooksList     []string        `json:"pendingBooksList"`
	CheckedOutBooksCount int64           `json:"checkedOutBooksCount"`
	CheckedOutBookList   []string        `json:"checkedOutBookList"`
	ReservedBooksCount   int64           `json:"reservedBooksCount"`
	ReservedBookList     []string        `json:"reservedBookList"`
	CompletedBooksCount  int64           `json:"completedBooksCount"`
	CompletedBooksList   []string        `json:"completedBooksList"`
	FavoriteGenres       []BookGenreType `json:"favoriteGenres"`
	WishlistBooks        []string        `json:"wishlistBooks"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            *time.Time      `json:"updatedAt"`
}

// reserved - books that are reseved but yet to be checked out
// checkedOut - books that are checked out
// pending - books that are checkedout but are yet to be added review or returned
// completed - books that are both added review and returned
