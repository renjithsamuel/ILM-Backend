package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"
)

// create checkout handler creates new checkout with given data
func (th *LibraryHandler) CreateCheckoutHandler(c *gin.Context) {
	req := model.CreateCheckoutRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// create checkout is an upsert operation
	if err := th.domain.CreateCheckoutTicket(&req); err != nil {
		if errors.Is(domain.ErrPaymentPending, err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "checkout ticket created successfully",
	})
}
