package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

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

	// encrypting the password of the user before stroign
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while hashing the password",
		})
		return
	}
	req.Password = string(hashedPassword)

	fmt.Println(req.Password)

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
