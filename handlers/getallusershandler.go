package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// get all users handler retrieves all users
func (th *LibraryHandler) GetAllUsersHandler(c *gin.Context) {
	// todo pagination
	// todo sort
	req := model.GetAllUsersRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
	}

	// Retrieve all users from the domain
	users, err := th.domain.GetAllUsers(&req)
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
