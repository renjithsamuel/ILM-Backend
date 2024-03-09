package domain

import (
	"errors"
	"integrated-library-service/model"

	"github.com/rs/zerolog/log"
)

var (
	// ErrFailedCreateBook is an error when create book failed
	ErrFailedCreateBook = errors.New("create book failed")
)

// create book creates new book
func (l *LibraryService) CreateBook(book *model.CreateBookRequest) error {
	sqlStatement := `
						INSERT INTO "users"(
									"profileImageUrl",
									"name",
									"email",
									"role"
								) VALUES ($1, $2, $3, $4)
					ON CONFLICT ("email") 
						DO UPDATE SET 
								"profileImageUrl" = EXCLUDED."profileImageUrl",
								"name" = EXCLUDED."name",
								"role" = EXCLUDED."role";
					`
	if err := l.db.QueryRow(sqlStatement, book); err != nil {
		log.Error().Msgf("[Error] CreateUser(), db.Exec err: %v", err)
		return ErrFailedCreateBook
	}

	return nil
}
