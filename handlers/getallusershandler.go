package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all users handler retrieves all users
func (th *LibraryHandler) GetAllUsersHandler(c *gin.Context) {
	// sort things

	// Retrieve all users from the domain
	users, err := th.domain.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the list of users in the response
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
