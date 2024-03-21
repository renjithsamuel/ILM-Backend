package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// GetAllBooksByBookDetailsFromHandler retrieves all books matching the condition [this is used in user book list]
func (th *LibraryHandler) GetAllBooksByBookDetailsFromHandler(c *gin.Context) {
	req := model.GetAllBooksByBookDetailsFromRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}
	// sort things need to be added

	// Retrieve all books from the domain
	books, err := th.domain.GetAllBooksByBookDetailsFrom(&model.GetAllBooksByBookDetailsFromRequest{
		UserID:          req.UserID,
		BookDetailsFrom: req.BookDetailsFrom,
	})
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
