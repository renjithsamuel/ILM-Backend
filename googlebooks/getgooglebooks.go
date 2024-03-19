package googlebooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"integrated-library-service/model"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// ErrGetGoogleBooks is an error when get books failed
	ErrGetGoogleBooks = errors.New("get google books failed")
)

// GetGoogleBooks gets required google books
func (g *GoogleBooksClient) GetGoogleBooks(request *model.GetAllNewBooksRequest) ([]*model.CreateBookRequest, int, error) {
	// requestBody, err := json.Marshal(&createCustomerAccountRequest)
	// if err != nil {
	// 	log.Error().Msgf("marshal error : %v ", err)
	// 	return nil, &CustomerAccountCreateResponse{Response: ErrorResponse(err, false)}
	// }
	startIndex := (request.Page - 1) * request.Limit
	url := fmt.Sprintf("/v1/volumes?q=orderBy=newest&startIndex=%v&maxResults=%v&key=%v", startIndex, request.Limit, g.apiKey)
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

func convertAndFormatDate(dateStr string) (*time.Time, error) {
	var layout string
	// Check if the length is 4 to handle YYYY format
	if len(dateStr) == 4 {
		layout = "2006"
	} else {
		// Try parsing with YYYY-MM-DD format first (common case)
		layout = "2006-01-02"
		_, err := time.Parse(layout, dateStr)
		// If parsing fails, try other layouts (add more as needed)
		if err != nil {
			// Example trying alternative layout (e.g., DD-MM-YYYY)
			layout = "02-01-2006"
			_, err = time.Parse(layout, dateStr)
			if err != nil {
				return nil, fmt.Errorf("invalid date format: %s", dateStr)
			}
		}
	}

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, err
	}
	time := t.UTC()
	return &time, nil
}
