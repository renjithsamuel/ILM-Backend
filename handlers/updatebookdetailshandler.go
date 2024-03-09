package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// UpdateUserHandler updates user details
func (th *LibraryHandler) UpdateBookDetailsHandler(c *gin.Context) {
	userIDInterface, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "userID not found",
		})
		return
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "userID is not of type string",
		})
		return
	}
	
	req := model.BookDetails{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
	}

	// update book details operation
	if err := th.domain.UpdateBookDetails(&req, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated book details successfully",
	})
}
