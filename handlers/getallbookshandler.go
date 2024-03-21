package handlers

import (
	"integrated-library-service/apperror"
	"integrated-library-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all books handler retrieves all books
func (th *LibraryHandler) GetAllBooksHandler(c *gin.Context) {
	req := model.GetAllBooksRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
	}

	// Retrieve all books from the domain
	books, err := th.domain.GetAllBooks(&req)
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
