package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
	"integrated-library-service/domain"
)

// get book by ISBN handler retrieves a book by its ISBN
func (th *LibraryHandler) GetBookByISBNHandler(c *gin.Context) {
	// Get the book ID from the URL parameter
	req := model.GetBookByISBNRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Retrieve the book from the domain
	book, err := th.domain.GetBookByISBN(req.ISBN)
	if err != nil {
		if err == domain.ErrGetBookByIDNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Book not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the book details in the response
	c.JSON(http.StatusOK, gin.H{
		"book": book,
	})
}
