package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardDataBoardHandler gets dashboard data
func (th *LibraryHandler) GetDashboardDataBoardHandler(c *gin.Context) {
	dashboardData, err := th.domain.GetDashboardDataBoard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboardData": dashboardData,
	})
}
