package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// get all books handler retrieves all books from a given specific list
func (th *LibraryHandler) GetAllNewBooksHandler(c *gin.Context) {
	// sort things need to be added
	req := model.GetAllBooksRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// get books from google
	googleBooks, totalPages, err := th.googleBooksService.GetGoogleBooks(&req)
	if err != nil {
		// temp solution
		// todo pagination
		books, err := th.domain.GetAllBooks(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"totalPages": -1,
			"books":      books,
		})
		return
		// c.JSON(http.StatusInternalServerError, gin.H{
		// 	"message": err.Error(),
		// })
		// return
	}

	// create books in batch
	if err := th.domain.CreateBooksBatch(googleBooks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	bookList := []string{}
	for _, book := range googleBooks {
		bookList = append(bookList, book.ISBN)
	}

	// Retrieve all books from the domain
	books, err := th.domain.GetAllBooksFromSpecific(bookList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Return the list of books in the response
	c.JSON(http.StatusOK, gin.H{
		"totalPages": totalPages,
		"books":      books,
	})
}
