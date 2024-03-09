package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// register user handler creates new user with given data
func (th *LibraryHandler) RegisterUserHandler(c *gin.Context) {
	req := model.RegisterUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// create user is an upsert operation
	if err := th.domain.CreateUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}
