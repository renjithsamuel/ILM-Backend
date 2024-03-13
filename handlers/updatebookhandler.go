package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"
)

// UpdateBookHandler updates an existing book
func (th *LibraryHandler) UpdateBookHandler(c *gin.Context) {
	var req model.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	if err := th.domain.UpdateBook(&req); err != nil {
		if err == domain.ErrUpdateBookNotFound {
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

	c.JSON(http.StatusOK, gin.H{
		"message": "book updated successfully",
	})
}
