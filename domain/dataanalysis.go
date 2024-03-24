package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"integrated-library-service/model"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	// ErrGetBooksByApproximateDemandFailed is when get approximate demand books failed
	ErrGetBooksByApproximateDemandFailed = errors.New("get approximate demand books failed")
)

// GetBooksByApproximateDemand
func (l *LibraryService) GetBooksByApproximateDemand(request *model.GetBooksByApproximateDemandRequest) ([]model.Book, uint, error) {
	sqlStatement := `
        SELECT 
			"ID",
			"ISBN",
			"title",
			"author",
			"genre",
			"publishedDate",
			"desc",
			"previewLink",
			"coverImage",
			"shelfNumber",
			"inLibrary",
			"views",
			"booksLeft",
			"wishlistCount",
			"rating",
			"reviewCount",
			"approximateDemand",
			"createdAt",
			"updatedAt",
			"reviewsList",
			"viewsList",
			"wishList"
        FROM 
            "books"
        WHERE
            "inLibrary" = true AND
            "rating" != 0 AND "views" != 0 AND "wishlistCount" != 0 AND  "reviewCount" != 0
        ORDER BY 
            ("rating" * $1) + ("views" * $2) + ("wishlistCount" * $3) + ("reviewCount" * $4) DESC -- orders by approximate demand points
        %s; -- limit and offset
    `

	// adding pagination to sqlstatement
	limitOffset := ` LIMIT %d OFFSET %d`
	limitOffset = fmt.Sprintf(limitOffset, request.Limit, (request.Page-1)*(request.Limit))
	sqlStatement = fmt.Sprintf(sqlStatement, limitOffset)

	rows, err := l.db.Query(sqlStatement, model.RatingWeight, model.ViewWeight, model.WishlistWeight, model.ReviewWeight)
	if err != nil {
		log.Error().Msgf("[Error] UpdateApproximateDemand(), db.Query err: %v", err)
		return nil, 0, nil
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var (
			book       model.Book
			updatedAt  sql.NullTime
			reviewList pq.StringArray
			viewList   pq.StringArray
			wishList   pq.StringArray
		)
		err := rows.Scan(
			&book.ID,
			&book.ISBN,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.PublishedDate,
			&book.Description,
			&book.PreviewLink,
			&book.CoverImage,
			&book.ShelfNumber,
			&book.InLibrary,
			&book.Views,
			&book.BooksLeft,
			&book.WishlistCount,
			&book.Rating,
			&book.ReviewCount,
			&book.ApproximateDemand,
			&book.CreatedAt,
			&updatedAt,
			&reviewList,
			&viewList,
			&wishList,
		)
		if err != nil {
			log.Error().Msgf("[Error] UpdateApproximateDemand(), rows.Scan err: %v", err)
			return nil, 0, nil
		}
		book.UpdatedAt = &updatedAt.Time
		book.ReviewsList = reviewList
		book.ViewsList = viewList
		book.WishList = wishList

		// get ratings from helper
		ratings, err := l.getAverageRating(book.ID)
		if err != nil && errors.Is(err, ErrRatingNotFound) {
			log.Error().Msgf("[Error] UpdateApproximateDemand(), getAverageRating err: %v", err)
			return nil, 0, nil
		}
		book.Rating = *ratings.Rating

		// Calculate demand score
		demandScore := l.calculateDemandScore(book)

		// Use demand score to estimate approximate demand
		approximateDemand := demandScore / model.RatingWeight // RatingWeight has the highest weight
		book.ApproximateDemand = approximateDemand

		books = append(books, book)
	}

	// total pages
	sqlStatementCount := `
	SELECT 
		COUNT(*)
	FROM 
		"books"
	WHERE
		"inLibrary" = true AND
		"rating" != 0 AND "views" != 0 AND "wishlistCount" != 0 AND  "reviewCount" != 0
	%s; --LIMIT AND OFFSET 	
`

	var totalRows uint
	err = l.db.QueryRow(fmt.Sprintf(sqlStatementCount, limitOffset)).Scan(&totalRows)
	// no rows
	if errors.Is(err, sql.ErrNoRows) {
		return []model.Book{}, 0, nil
	}
	if err != nil {
		log.Error().Msgf("[Error] UpdateApproximateDemand(), count query err: %v", err)
		return nil, 0, err
	}
	// Calculate total pages
	totalPages := (uint32(totalRows) + request.Limit - 1) / request.Limit

	return books, uint(totalPages), nil
}
