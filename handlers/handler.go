package handlers

import (
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
)

var (
	validate = validator.New()
)

func init() {
	apperror.RegisterTags(validate)
}

// Handler is an interface for library operations.
type Handler interface {
	HealthHandler(c *gin.Context)
	// user related
	RegisterUserHandler(c *gin.Context)
	LoginUserHandler(c *gin.Context)
	GetUserHandler(c *gin.Context)
	GetAllUsersHandler(c *gin.Context)
	UpdateUserHandler(c *gin.Context)
	UpdateBookDetailsHandler(c *gin.Context)
	DeleteUserHandler(c *gin.Context)
	// book related
	CreateBookHandler(c *gin.Context)
	GetBookByISBNHandler(c *gin.Context)
	GetAllBooksHandler(c *gin.Context)
	// checkout related
	
}

type LibraryHandler struct {
	domain    domain.Service
	secretKey string
}

// NewLibraryHandler returns new instance of Handler.
func NewLibraryHandler(domain domain.Service, secretKey string) *LibraryHandler {
	h := &LibraryHandler{
		domain:    domain,
		secretKey: secretKey,
	}

	return h
}
