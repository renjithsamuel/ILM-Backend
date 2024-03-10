package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"
)

// UpdateCheckoutTicketHandler updates an existing checkout ticket
func (th *LibraryHandler) UpdateCheckoutTicketHandler(c *gin.Context) {
	checkoutID := c.Param("id")

	req := model.UpdateCheckoutTicketRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Check if the checkout ticket exists
	_, err := th.domain.GetCheckoutTicketByID(checkoutID)
	if err != nil {
		if errors.Is(err, domain.ErrGetCheckoutTicketByIDNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		return
	}

	// Update the checkout ticket using the domain function
	err = th.domain.UpdateCheckoutTicket(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update checkout ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "checkout ticket updated successfully",
	})
}
