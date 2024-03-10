package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllCheckoutTicketsHandler retrieves all checkout tickets
func (th *LibraryHandler) GetAllCheckoutTicketsHandler(c *gin.Context) {
	// get sorting data

	// Retrieve all checkout tickets using the domain function
	checkoutTickets, err := th.domain.GetAllCheckoutTickets()
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
