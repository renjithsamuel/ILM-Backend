package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all books handler retrieves all books
func (th *LibraryHandler) GetAllBooksHandler(c *gin.Context) {
	// sort things need to be added
	
	// Retrieve all books from the domain
	books, err := th.domain.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the list of books in the response
	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}
