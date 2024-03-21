package handlers

import (
	"integrated-library-service/apperror"
	"integrated-library-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetReviewsByBookIDHandler retrieves all reviews for a particular book
func (th *LibraryHandler) GetReviewsByBookIDHandler(c *gin.Context) {
	// todo pagination
	req := struct {
		BookID string `json:"bookID" uri:"bookid" binding:"required,uuid"`
	}{}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	reqPagination := model.ReviewSort{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
	}

	reviews, err := th.domain.GetReviewsByBookID(req.BookID, &reqPagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
	})
}
