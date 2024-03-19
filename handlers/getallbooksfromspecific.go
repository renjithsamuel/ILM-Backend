package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
)

// get all books handler retrieves all books from a given specific list
func (th *LibraryHandler) GetAllBooksFromSpecificHandler(c *gin.Context) {
	// sort things need to be added
	req := struct {
		BooksList []string `json:"booksList" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&req.BooksList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Retrieve all books from the domain
	books, err := th.domain.GetAllBooksFromSpecific(req.BooksList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the list of books in the response
	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}
