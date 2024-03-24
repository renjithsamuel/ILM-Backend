package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"
)

// GetRecommendedBooksForUserHandler gets paginated recommneded books for user
func (th *LibraryHandler) GetRecommendedBooksForUserHandler(c *gin.Context) {
	userIDInterface, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "userID not found",
		})
		return
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "userID is not of type string",
		})
		return
	}

	req := model.GetRecommendedBooksForUserRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	user, err := th.domain.GetUserWithBookDetails(userID)
	if err != nil {
		if errors.Is(err, domain.ErrGetUserWithBookDetailsNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	bookList := th.getBookListFromUserBookDetails(user.BookDetails)

	// Retrieve all books from the domain
	books, err := th.domain.GetAllBooksFromSpecific(bookList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// generate search text from books
	searchText := th.generateSearchTextFromBooks(books)

	searchRequest := &model.SearchRequest{
		Page:       req.Page,
		Limit:      req.Limit,
		SortBy:     "title",
		OrderBy:    "ascending",
		SearchBy:   "default",
		Type:       "book",
		SearchText: searchText,
	}

	// get books from google
	googleBooks, totalPages, err := th.googleBooksService.SearchGoogleBooks(searchRequest)
	if err != nil {
		// temp solution
		books, totalPages, err := th.domain.GetAllBooksForSearch(searchRequest)
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

	bookList = []string{}
	for _, book := range googleBooks {
		bookList = append(bookList, book.ISBN)
	}

	// Retrieve all books from the domain
	books, err = th.domain.GetAllBooksFromSpecific(bookList)
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

// getBookListFromUserBookDetails
func (th *LibraryHandler) getBookListFromUserBookDetails(bookDetails model.BookDetails) []string {
	bookList := []string{}
	bookList = append(bookList, bookDetails.WishlistBooks...)
	bookList = append(bookList, bookDetails.CompletedBooksList...)
	bookList = append(bookList, bookDetails.CheckedOutBookList...)
	bookList = append(bookList, bookDetails.ReservedBookList...)
	return bookList
}

// generateSearchTextFromBooks
func (th *LibraryHandler) generateSearchTextFromBooks(books []model.Book) string {
	var searchText string

	for _, book := range books {
		searchText += book.Genre + book.Title + book.Author
	}

	return searchText
}
