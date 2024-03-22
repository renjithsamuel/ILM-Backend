package domain

import (
	"database/sql"

	"integrated-library-service/model"
)

// Service is an interface for its concrete class to implement.
type Service interface {
	DBStatus() (bool, error)
	// user related
	CreateUser(user *model.RegisterUserRequest) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserWithBookDetails(userID string) (*model.User, error)
	GetAllUsers(request *model.GetAllUsersRequest) ([]model.User, uint32, error)
	GetAllUsersForSearch(request *model.SearchRequest) ([]model.User, uint, error)
	UpdateUser(user *model.User, userID string) error
	UpdateBookDetails(bookDetails *model.BookDetails, userID string) error
	DeleteUser(userID string) error
	// book related
	CreateBook(book *model.CreateBookRequest) error
	GetBookByISBN(ISBN string) (*model.Book, error)
	GetAllBooks(request *model.GetAllBooksRequest) ([]model.Book, uint, error)
	GetAllBooksForSearch(request *model.SearchRequest) ([]model.Book, uint, error)
	GetAllBooksByBookDetailsFrom(request *model.GetAllBooksByBookDetailsFromRequest) ([]model.Book, error)
	GetAllBooksFromSpecific(request []string) ([]model.Book, error)
	CreateBooksBatch(books []*model.CreateBookRequest) error
	UpdateBook(book *model.UpdateBookRequest) error
	// checkout related
	CreateCheckoutTicket(ticket *model.CreateCheckoutRequest) error
	GetCheckoutTicketByID(ticketID string) (*model.CheckoutTicket, error)
	GetCheckoutsByUserID(bookID, userID string) ([]model.CheckoutTicket, error)
	GetAllCheckoutTicketsWithDetails(request *model.GetAllCheckoutData) ([]model.CheckoutTicketResponse, uint, error)
	UpdateCheckoutTicket(ticket *model.UpdateCheckoutTicketRequest) error
	DeleteCheckoutTicket(ticketID string) error
	// review related
	CreateReview(review *model.CreateReviewRequest) error
	GetReviewByID(reviewID string) (*model.Review, error)
	UpdateReview(updateReq *model.UpdateReviewRequest) error
	DeleteReview(reviewID string) error
	GetReviewsByBookID(bookID string, sortPagination *model.ReviewSort) ([]model.Review, uint, error)
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
