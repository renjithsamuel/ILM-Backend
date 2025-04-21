package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
)

// DeleteReviewHandler deletes a review by its ID
func (th *LibraryHandler) DeleteReviewHandler(c *gin.Context) {
	req := struct {
		ReviewID string `json:"reviewID" uri:"reviewid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	if err := th.domain.DeleteReview(req.ReviewID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review deleted successfully",
	})
}
