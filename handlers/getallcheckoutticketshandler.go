package handlers

import (
	"integrated-library-service/apperror"
	"integrated-library-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllCheckoutTicketsHandler retrieves all checkout tickets
func (th *LibraryHandler) GetAllCheckoutTicketsHandler(c *gin.Context) {
	req := model.GetAllCheckoutData{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
	}
	// Retrieve all checkout tickets using the domain function
	checkoutTickets, err := th.domain.GetAllCheckoutTicketsWithDetails(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkoutTickets": checkoutTickets,
	})
}
