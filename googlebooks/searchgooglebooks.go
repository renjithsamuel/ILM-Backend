package googlebooks

import (
	"encoding/json"
	"fmt"
	"integrated-library-service/model"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

// GetGoogleBooks gets required google books
func (g *GoogleBooksClient) SearchGoogleBooks(request *model.SearchRequest) ([]*model.CreateBookRequest, int, error) {
	startIndex := ((request.Page) - 1) * (request.Limit)

	searchQuery := request.SearchText

	switch request.SearchBy {
	case "title":
		searchQuery = fmt.Sprintf("intitle:%s", request.SearchText)
	case "author":
		searchQuery = fmt.Sprintf("inauthor:%s", request.SearchText)
	case "isbn":
		searchQuery = fmt.Sprintf("isbn:%s", request.SearchText)
	case "subject": // used for similar books and nowhere else
		searchQuery = fmt.Sprintf("subject:%s", request.SearchText)
	}

	// Encode the search query to replace spaces with %20
	encodedSearchQuery := url.QueryEscape(searchQuery)

	// Construct the URL
	url := fmt.Sprintf("/v1/volumes?q=%s&orderBy=%v&startIndex=%v&maxResults=%v&key=%v", encodedSearchQuery, "relevance", startIndex, request.Limit, g.apiKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Msgf("http request error : %v ", err)
		return nil, 0, ErrGetGoogleBooks
	}

	resp, err := g.do(req)
	if err != nil {
		log.Error().Msgf("response error : %v ", err)
		return nil, 0, ErrGetGoogleBooks
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Error().Msgf("getGoogleBooks error : %v , status : %v ", bodyBytes, resp.StatusCode)
		return nil, 0, ErrGetGoogleBooks
	}

	var googleBooks model.GoogleBookResponse
	if err = json.NewDecoder(resp.Body).Decode(&googleBooks); err != nil {
		log.Error().Msgf("decode body error : %v ", err)
		return nil, 0, ErrGetGoogleBooks
	}

	books := []*model.CreateBookRequest{}
	for _, googleBook := range googleBooks.Items {
		book := model.CreateBookRequest{
			Title:       googleBook.VolumeInfo.Title,
			CoverImage:  googleBook.VolumeInfo.ImageLinks.Thumbnail,
			Description: googleBook.VolumeInfo.Subtitle,
			PreviewLink: googleBook.VolumeInfo.PreviewLink,
		}
		if len(googleBook.VolumeInfo.IndustryIdentifiers) > 0 {
			book.ISBN = googleBook.VolumeInfo.IndustryIdentifiers[0].Identifier
		} else {
			continue
		}

		if len(googleBook.VolumeInfo.Categories) > 0 {
			book.Genre = googleBook.VolumeInfo.Categories[0]
		} else {
			book.Genre = "other"
		}

		if len(googleBook.VolumeInfo.Authors) > 0 {
			book.Author = googleBook.VolumeInfo.Authors[0]
		} else {
			book.Author = "unknown"
		}

		t, err := convertAndFormatDate(googleBook.VolumeInfo.PublishedDate)
		if err != nil {
			log.Error().Msgf("Error while converting string to data %v", err)
		} else {
			book.PublishedDate = *t
		}

		// declaring default values
		mockvalue := int64(0)
		mockFloat := float64(0)
		mockBool := false
		book.BooksLeft = &mockvalue
		book.ShelfNumber = &mockvalue
		book.InLibrary = &mockBool
		book.Rating = &mockFloat
		book.ReviewCount = &mockvalue
		book.Views = &mockvalue
		book.ReviewsList = []string{}
		book.ViewsList = []string{}
		book.WishList = []string{}
		book.WishlistCount = &mockvalue
		book.ApproximateDemand = &mockvalue

		books = append(books, &book)
	}

	return books, googleBooks.TotalItems, nil
}
