package domain

import (
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"integrated-library-service/model"
)

var (
	// ErrReviewNotFound is an error when get ratings not found
	ErrRatingNotFound = errors.New("get ratings not found")
	// ErrReviewNotFound is an error when get ratings not found
	ErrGetRating = errors.New("get ratings failed")
)

// getAverageRating fetchs overall product specifics
func (r *LibraryService) getAverageRating(bookID string) (*model.GetAverageRatingResponse, error) {
	sqlSt := `	SELECT "rating", COUNT("rating")
				FROM reviews
				WHERE "bookID" = $1
				GROUP BY "rating"
			`
	ratingEntity := model.RatingEntity{}
	rows, err := r.db.Query(sqlSt, bookID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			log.Error().Msgf("FetchReviews: DB fetch query error: %v", err)
			return &model.GetAverageRatingResponse{
				RatingEntity: ratingEntity,
			}, ErrRatingNotFound
		}
		log.Error().Msgf("FetchReviews: DB fetch query error: %v", err)
		return &model.GetAverageRatingResponse{
			RatingEntity: ratingEntity,
		}, ErrGetRating
	}
	defer rows.Close()

	var totalReviews uint32 = 0
	var productRating float64

	for rows.Next() {
		var rating uint8
		var count uint32
		err := rows.Scan(&rating, &count)
		if err != nil {
			log.Error().Msgf("FetchReviews: DB row scan error: %v", err)
			return &model.GetAverageRatingResponse{
				RatingEntity: ratingEntity,
			}, err
		}
		totalReviews += count
		switch rating {
		case 1:
			ratingEntity.OneStar += count
		case 2:
			ratingEntity.TwoStar += count
		case 3:
			ratingEntity.ThreeStar += count
		case 4:
			ratingEntity.FourStar += count
		case 5:
			ratingEntity.FiveStar += count
		default:
			log.Error().Msgf("rating value %v is not supported", rating)
		}
	}

	totalRating := ratingEntity.OneStar*1 + ratingEntity.TwoStar*2 + ratingEntity.ThreeStar*3 +
		ratingEntity.FourStar*4 + ratingEntity.FiveStar*5

	if totalRating != 0 {
		productRating = (float64(totalRating) / float64(totalReviews)) * 100 / 100
	}

	return &model.GetAverageRatingResponse{Rating: &productRating, TotalReviews: totalReviews,
		RatingEntity: ratingEntity,
	}, nil
}

// Function to calculate the demand score for a book
func (l *LibraryService) calculateDemandScore(book model.Book) int64 {
	// Calculate points for each factor
	ratingPoints := int(book.Rating * model.RatingWeight)
	reviewPoints := len(book.ReviewsList) * model.ReviewWeight
	viewPoints := int(book.Views * model.ViewWeight)
	wishlistPoints := int(book.WishlistCount * model.WishlistWeight)

	// Calculate cumulative demand score
	demandScore := ratingPoints + reviewPoints + viewPoints + wishlistPoints
	return int64(demandScore)
}
