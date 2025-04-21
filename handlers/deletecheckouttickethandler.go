package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
)

// DeleteCheckoutTicketHandler deletes a checkout ticket by its ID
func (th *LibraryHandler) DeleteCheckoutTicketHandler(c *gin.Context) {
	req := struct {
		CheckoutID string `json:"checkoutid" uri:"checkoutid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Delete the checkout ticket using the domain function
	err := th.domain.DeleteCheckoutTicket(req.CheckoutID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "checkout ticket deleted successfully",
	})
}
