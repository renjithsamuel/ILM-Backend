package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/model"
)

// search handler returns either user data or book data from DB or google books respectively
func (th *LibraryHandler) SearchHandler(c *gin.Context) {
	// sort things need to be added
	req := model.SearchRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	if *req.Type == model.SearchRequestTypeBook {
		// get books from google
		googleBooks, totalPages, err := th.googleBooksService.SearchGoogleBooks(&req)
		if err != nil {
			books, totalPages, err := th.domain.GetAllBooksForSearch(&req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"totalPages": totalPages,
				"books":      books,
			})
			return
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

	if *req.Type == model.SearchRequestTypeUser {
		// Retrieve all books from the domain
		users, totalPages, err := th.domain.GetAllUsersForSearch(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Return the list of books in the response
		c.JSON(http.StatusOK, gin.H{
			"totalPages": totalPages,
			"users":      users,
		})
	}

}
