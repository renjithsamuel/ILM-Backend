package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
)

// GetReviewByIDHandler retrieves a review by its ID
func (th *LibraryHandler) GetReviewByIDHandler(c *gin.Context) {
	req := struct {
		ReviewID string `json:"reviewID" uri:"reviewid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	review, err := th.domain.GetReviewByID(req.ReviewID)
	if err != nil {
		if errors.Is(domain.ErrGetReviewByIDNotFound, err) {
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
		"review": review,
	})
}
