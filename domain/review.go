package domain

import (
	"database/sql"
	"errors"
	"time"

	"integrated-library-service/model"

	"github.com/rs/zerolog/log"
)

var (
	// ErrFailedCreateReview is an error when create review failed
	ErrFailedCreateReview = errors.New("create review failed")
	// ErrGetReviewByIDFailed is an error when get review failed
	ErrGetReviewByIDFailed = errors.New("get review failed")
	// ErrGetReviewsFailed is an error when get reviews failed
	ErrGetReviewsFailed = errors.New("get reviews failed")
	// ErrGetReviewByIDNotFound is an error when get review not found
	ErrGetReviewByIDNotFound = errors.New("get review not found")
	// ErrFailedUpdateReview is an error when update review failed
	ErrFailedUpdateReview = errors.New("update review failed")
	// ErrFailedDeleteReview is an error when delete review failed
	ErrFailedDeleteReview = errors.New("delete review failed")
)

// CreateReview creates a new review
func (l *LibraryService) CreateReview(review *model.CreateReviewRequest) error {
	sqlStatement := `
		INSERT INTO "reviews"(
			"bookID",
			"checkoutID",
			"userID",
			"commentHeading",
			"comment",
			"rating",
			"createdAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) 
		ON CONFLICT ("bookID", "checkoutID", "userID") 
		DO UPDATE SET
			"commentHeading" = EXCLUDED."commentHeading",
			"comment" = EXCLUDED."comment",
			"rating" = EXCLUDED."rating",
			"updatedAt" = NOW()
		RETURNING "ID";
	`

	var reviewID string
	err := l.db.QueryRow(
		sqlStatement,
		review.BookID,
		review.CheckoutID,
		review.UserID,
		review.CommentHeading,
		review.Comment,
		review.Rating,
		time.Now().UTC(),
	).Scan(&reviewID)

	if err != nil {
		log.Error().Msgf("[Error] CreateReview(), db.QueryRow err: %v", err)
		return ErrFailedCreateReview
	}

	return nil
}

// GetReviewByID retrieves a review by its ID
func (l *LibraryService) GetReviewByID(reviewID string) (*model.Review, error) {
	sqlStatement := `
		SELECT 
			"ID",
			"bookID",
			"checkoutID",
			"userID",
			"commentHeading",
			"comment",
			"rating",
			"likes",
			"createdAt",
			"updatedAt"
		FROM 
			"reviews"
		WHERE 
			"ID" = $1;
	`

	var (
		review    model.Review
		updatedAt sql.NullTime
	)
	err := l.db.QueryRow(sqlStatement, reviewID).Scan(
		&review.ID,
		&review.BookID,
		&review.CheckoutID,
		&review.UserID,
		&review.CommentHeading,
		&review.Comment,
		&review.Rating,
		&review.Likes,
		&review.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msgf("[Error] GetReviewByID(), db.QueryRow err: %v", err)
			return nil, ErrGetReviewByIDNotFound
		}

		log.Error().Msgf("[Error] GetReviewByID(), db.QueryRow err: %v", err)
		return nil, ErrGetReviewByIDFailed
	}
	review.UpdatedAt = &updatedAt.Time

	return &review, nil
}


// GetReviewsByBookID retrieves all reviews for a particular book
func (l *LibraryService) GetReviewsByBookID(bookID string) ([]model.Review, error) {
	sqlStatement := `
		SELECT 
			"ID",
			"bookID",
			"checkoutID",
			"userID",
			"commentHeading",
			"comment",
			"rating",
			"likes",
			"createdAt",
			"updatedAt"
		FROM 
			"reviews"
		WHERE 
			"bookID" = $1;
	`

	rows, err := l.db.Query(sqlStatement, bookID)
	if err != nil {
		log.Error().Msgf("[Error] GetReviewsByBookID(), db.Query err: %v", err)
		return nil, ErrGetReviewsFailed
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var (
			review    model.Review
			updatedAt sql.NullTime
		)
		err := rows.Scan(
			&review.ID,
			&review.BookID,
			&review.CheckoutID,
			&review.UserID,
			&review.CommentHeading,
			&review.Comment,
			&review.Rating,
			&review.Likes,
			&review.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Error().Msgf("[Error] GetReviewsByBookID(), rows.Scan err: %v", err)
			return nil, ErrGetReviewsFailed
		}
		review.UpdatedAt = &updatedAt.Time
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// GetAllReviews retrieves all reviews
func (l *LibraryService) GetAllReviews() ([]model.Review, error) {
	sqlStatement := `
		SELECT 
			"ID",
			"bookID",
			"checkoutID",
			"userID",
			"commentHeading",
			"comment",
			"rating",
			"likes",
			"createdAt",
			"updatedAt"
		FROM 
			"reviews";
	`

	rows, err := l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetAllReviews(), db.Query err: %v", err)
		return nil, ErrGetReviewsFailed
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var (
			review    model.Review
			updatedAt sql.NullTime
		)
		err := rows.Scan(
			&review.ID,
			&review.BookID,
			&review.CheckoutID,
			&review.UserID,
			&review.CommentHeading,
			&review.Comment,
			&review.Rating,
			&review.Likes,
			&review.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			log.Error().Msgf("[Error] GetAllReviews(), rows.Scan err: %v", err)
			return nil, ErrGetReviewsFailed
		}
		review.UpdatedAt = &updatedAt.Time
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// UpdateReview updates an existing review
func (l *LibraryService) UpdateReview(review *model.UpdateReviewRequest) error {
	sqlStatement := `
		UPDATE "reviews" SET
			"commentHeading" = $2,
			"comment" = $3,
			"rating" = $4,
			"likes" = $5,
			"updatedAt" = $6
		WHERE
			"ID" = $1;
	`

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	_, err := l.db.Exec(
		sqlStatement,
		review.ID,
		review.CommentHeading,
		review.Comment,
		review.Rating,
		review.Likes,
		updatedAt,
	)

	if err != nil {
		log.Error().Msgf("[Error] UpdateReview(), db.Exec err: %v", err)
		return ErrFailedUpdateReview
	}

	return nil
}

// DeleteReview deletes a review by its ID
func (l *LibraryService) DeleteReview(reviewID string) error {
	sqlStatement := `
		DELETE FROM "reviews" WHERE "ID" = $1;
	`

	_, err := l.db.Exec(sqlStatement, reviewID)
	if err != nil {
		log.Error().Msgf("[Error] DeleteReview(), db.Exec err: %v", err)
		return ErrFailedDeleteReview
	}

	return nil
}
