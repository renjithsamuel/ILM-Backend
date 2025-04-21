package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// EmptyHandler is a handler that sends an OK status when called
func (th *LibraryHandler) EmptyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "passed",
	})
}
