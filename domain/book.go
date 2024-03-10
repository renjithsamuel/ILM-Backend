package domain

import (
	"database/sql"
	"errors"
	"integrated-library-service/model"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// ErrFailedCreateBook is an error when create book failed
	ErrFailedCreateBook = errors.New("create book failed")
	// ErrFailedGetBookByID is an error when get book failed
	ErrFailedGetBookByID = errors.New("get book failed")
	// ErrGetBookByIDNotFound is an error when get not found
	ErrGetBookByIDNotFound = errors.New("get book not found")
	// ErrFailedUpdateBook is an error when update book failed
	ErrFailedUpdateBook = errors.New("update book failed")
	// ErrFailedDeleteBook is an error when delete book failed
	ErrFailedDeleteBook = errors.New("delete book failed")
)

// CreateBook creates a new book or updates an existing one based on ISBN
func (l *LibraryService) CreateBook(book *model.Book) error {
	sqlStatement := `
		INSERT INTO "books"(
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
			"createdAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) 
		ON CONFLICT("ISBN") 
		DO UPDATE SET
			"title" = EXCLUDED."title",
			"author" = EXCLUDED."author",
			"genre" = EXCLUDED."genre",
			"publishedDate" = EXCLUDED."publishedDate",
			"desc" = EXCLUDED."desc",
			"previewLink" = EXCLUDED."previewLink",
			"coverImage" = EXCLUDED."coverImage",
			"shelfNumber" = EXCLUDED."shelfNumber",
			"inLibrary" = EXCLUDED."inLibrary",
			"views" = EXCLUDED."views",
			"booksLeft" = EXCLUDED."booksLeft",
			"wishlistCount" = EXCLUDED."wishlistCount",
			"rating" = EXCLUDED."rating",
			"reviewCount" = EXCLUDED."reviewCount",
			"approximateDemand" = EXCLUDED."approximateDemand",
			"updatedAt" = NOW()
		RETURNING "ID";
	`

	var bookID string
	err := l.db.QueryRow(
		sqlStatement,
		book.ISBN,
		book.Title,
		book.Author,
		book.Genre,
		book.PublishedDate,
		book.Description,
		book.PreviewLink,
		book.CoverImage,
		book.ShelfNumber,
		book.InLibrary,
		book.Views,
		book.BooksLeft,
		book.WishlistCount,
		book.Rating,
		book.ReviewCount,
		book.ApproximateDemand,
		book.CreatedAt,
	).Scan(&bookID)

	if err != nil {
		log.Error().Msgf("[Error] CreateBook(), db.QueryRow err: %v", err)
		return ErrFailedCreateBook
	}

	return nil
}

// GetBookByID retrieves a book by its ID
func (l *LibraryService) GetBookByISBN(ISBN string) (*model.Book, error) {
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
			"updatedAt"
		FROM
			"books"
		WHERE
			"ISBN" = $1;
	`

	var book model.Book
	err := l.db.QueryRow(sqlStatement, ISBN).Scan(
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
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			log.Error().Msgf("[Error] GetBookByID(), db.QueryRow err: %v", err)
			return nil, ErrGetBookByIDNotFound
		}
		log.Error().Msgf("[Error] GetBookByID(), db.QueryRow err: %v", err)
		return nil, ErrFailedGetBookByID
	}

	return &book, nil
}


// getAllBooks retrieves all books from the database
func (l *LibraryService) GetAllBooks() ([]model.Book, error) {
	sqlStatement := `
		SELECT 
			"ID",
			"ISBN",
			"title",
			"author",
			"genre",
			"publishedDate",
			"description",
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
			"updatedAt"
		FROM 
			"books";
	`

	rows, err := l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetAllBooks(), db.Query err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
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
			&book.UpdatedAt,
		)
		if err != nil {
			log.Error().Msgf("[Error] GetAllBooks(), rows.Scan err: %v", err)
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

// UpdateBook updates an existing book in the "books" table
func (l *LibraryService) UpdateBook(book *model.Book, bookID string) error {
	sqlStatement := `
		UPDATE "books" SET
			"ISBN" = $1,
			"title" = $2,
			"author" = $3,
			"genre" = $4,
			"publishedDate" = $5,
			"desc" = $6,
			"previewLink" = $7,
			"coverImage" = $8,
			"shelfNumber" = $9,
			"inLibrary" = $10,
			"views" = $11,
			"booksLeft" = $12,
			"wishlistCount" = $13,
			"rating" = $14,
			"reviewCount" = $15,
			"approximateDemand" = $16,
			"updatedAt" = $17
		WHERE
			"ID" = $18;
	`

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	res, err := l.db.Exec(
		sqlStatement,
		book.ISBN,
		book.Title,
		book.Author,
		book.Genre,
		book.PublishedDate,
		book.Description,
		book.PreviewLink,
		book.CoverImage,
		book.ShelfNumber,
		book.InLibrary,
		book.Views,
		book.BooksLeft,
		book.WishlistCount,
		book.Rating,
		book.ReviewCount,
		book.ApproximateDemand,
		updatedAt,
		bookID,
	)

	if err != nil {
		log.Error().Msgf("[Error] UpdateBook(), db.QueryRow err: %v", err)
		return ErrFailedUpdateBook
	}

	if rowsEffected, err := res.RowsAffected(); err != nil || rowsEffected == 0 {
		log.Error().Msgf("[error] UpdateBook(), [No rows affected]  : %v", err)
		return ErrFailedUpdateBook
	}

	return nil
}

// DeleteBook deletes a book from the "books" table based on bookID
func (l *LibraryService) DeleteBook(bookID string) error {
	sqlStatement := `
		DELETE FROM "books" WHERE "ID" = $1;
	`

	_, err := l.db.Exec(sqlStatement, bookID)
	if err != nil {
		log.Error().Msgf("[Error] DeleteBook(), db.Exec err: %v", err)
		return ErrFailedDeleteBook
	}

	return nil
}
