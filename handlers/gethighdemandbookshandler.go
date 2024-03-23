package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHighDemandBooksHandler gets high demand books
func (th *LibraryHandler) GetHighDemandBooksHandler(c *gin.Context) {
	highDemandBooks, err := th.domain.GetHighDemandBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"highDemandBooks": highDemandBooks,
	})
}
