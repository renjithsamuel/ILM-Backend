package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// GetApproximateDemandHandler return paginated books sorted based on approximate demand books
func (th *LibraryHandler) GetApproximateDemandHandler(c *gin.Context) {
	req := model.GetBooksByApproximateDemandRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Retrieve all books from the domain
	books, totalPages, err := th.domain.GetBooksByApproximateDemand(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the list of books in the response
	c.JSON(http.StatusOK, gin.H{
		"totalPages": totalPages,
		"books":      books,
	})
}
