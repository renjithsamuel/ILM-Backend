package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// delete user handler deletes user with book details
func (th *LibraryHandler) DeleteUserHandler(c *gin.Context) {
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

	// delete user
	if err := th.domain.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
