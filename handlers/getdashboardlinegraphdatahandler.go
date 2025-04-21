package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardLineGraphDataHandler gets line graph data
func (th *LibraryHandler) GetDashboardLineGraphDataHandler(c *gin.Context) {
	graphData, err := th.domain.GetDashboardLineGraphData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"graphData": graphData,
	})
}
