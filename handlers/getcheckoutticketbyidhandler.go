package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
)

// GetCheckoutTicketByIDHandler retrieves a checkout ticket by its ID
func (th *LibraryHandler) GetCheckoutTicketByIDHandler(c *gin.Context) {
	req := struct {
		CheckoutID string `json:"checkoutid" uri:"checkoutid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Retrieve the checkout ticket using the domain function
	checkoutTicket, err := th.domain.GetCheckoutTicketByID(req.CheckoutID)
	if err != nil {
		if errors.Is(domain.ErrGetCheckoutTicketByIDNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkoutTicket": checkoutTicket,
	})
}
