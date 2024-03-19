package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"
)

// GetCheckoutByUserIDHandler retrieves a checkout ticket by its UserID and BookID
func (th *LibraryHandler) GetCheckoutsByUserIDHandler(c *gin.Context) {
	req := struct {
		BookID string `json:"bookid" uri:"bookid" binding:"required,uuid"`
		UserID string `json:"userid" uri:"userid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
	}

	// Retrieve the checkout ticket using the domain function
	checkoutTickets, err := th.domain.GetCheckoutsByUserID(req.BookID, req.UserID)
	if err != nil {
		if errors.Is(domain.ErrGetCheckoutByUserIDNotFound, err) {
			c.JSON(http.StatusOK, gin.H{
				"checkoutTickets": []model.CheckoutTicket{},
			})
			return
		}	
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkoutTickets": checkoutTickets,
	})
}
