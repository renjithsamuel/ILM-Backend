package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"
)

// SimilarBooksHandler returns similar book to the provided book ISBN
func (th *LibraryHandler) SimilarBooksHandler(c *gin.Context) {
	req := model.SimilarBooksRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	book, err := th.domain.GetBookByISBN(req.ISBN)
	if err != nil {
		if err == domain.ErrGetBookByIDNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Book not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	searchRequest := &model.SearchRequest{
		Page:       1,
		Limit:      3,
		SortBy:     "title",
		OrderBy:    "ascending",
		SearchBy:   "recommendation",
		Type:       "book",
		SearchText: book.Genre + book.Author + book.Description,
	}

	// get similar books form google
	googleBooks, _, err := th.googleBooksService.SearchGoogleBooks(searchRequest)
	if err != nil {
		// get similar books with domain
		books, _, err := th.domain.GetAllBooksForSearch(searchRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		// if 3 books are not there return any 3 book
		if len(books) < 3 {
			searchRequestTemp := searchRequest
			searchRequestTemp.SearchText = ""
			books, _, err := th.domain.GetAllBooksForSearch(searchRequest)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"books": books,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"books": books,
		})
		return
	}

	// if 3 books are not there return any 3 book
	if len(googleBooks) < 3 {
		searchRequestTemp := searchRequest
		searchRequestTemp.SearchText = ""
		books, _, err := th.domain.GetAllBooksForSearch(searchRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"books": books,
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
		"books": books,
	})

}
