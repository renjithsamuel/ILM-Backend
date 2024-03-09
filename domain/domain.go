package domain

import (
	"database/sql"

	"integrated-library-service/model"
)

// Service is an interface for its concrete class to implement.
type Service interface {
	DBStatus() (bool, error)
	CreateUser(user *model.RegisterUserRequest) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserWithBookDetails(userID string) (*model.User, error)
	UpdateUser(user *model.User, userID string) error
	UpdateBookDetails(bookDetails *model.BookDetails, userID string) error
	DeleteUser(userID string) error
}

// LibraryService is a concrete service which implements Service
type LibraryService struct {
	db *sql.DB
}

// NewLibraryService is a constructor which creates an object of the LibraryService class.
func NewLibraryService(db *sql.DB) *LibraryService {
	return &LibraryService{
		db: db,
	}
}
