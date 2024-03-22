package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
)

// GetUserByIDHandler gets user with user with bookdetails for the given userid
func (th *LibraryHandler) GetUserByIDHandler(c *gin.Context) {
	req := struct {
		UserID string `json:"userid" uri:"userid" binding:"required,uuid"`
	}{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}
	// get user
	user, err := th.domain.GetUserWithBookDetails(req.UserID)
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
