package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/domain"
)

// get user handler gets user with book details
func (th *LibraryHandler) GetUserHandler(c *gin.Context) {
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

	// create user is an upsert operation
	user, err := th.domain.GetUserWithBookDetails(userID)
	if err != nil {
		if errors.Is(err, domain.ErrGetUserWithBookDetailsNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
