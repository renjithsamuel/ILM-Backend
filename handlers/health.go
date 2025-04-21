package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns alive status.
func (th *LibraryHandler) HealthHandler(c *gin.Context) {
	ok, err := th.domain.DBStatus()
	if err != nil {
		c.JSON(http.StatusFailedDependency, gin.H{
			"status": "alive",
			"db":     err.Error(),
		})
		return
	}
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
			"db":     "connected",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "alive",
		"db":     "false",
	})
}
