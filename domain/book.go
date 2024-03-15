package domain

import (
	"database/sql"
	"errors"
	"integrated-library-service/model"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	// ErrFailedCreateBook is an error when create book failed
	ErrFailedCreateBook = errors.New("create book failed")
	// ErrFailedGetBookByID is an error when get book failed
	ErrFailedGetBookByID = errors.New("get book failed")
	// ErrGetAllBooksByBookDetailsFromFailed is error when getting books by book details from
	ErrGetAllBooksByBookDetailsFromFailed = errors.New("get books from book details failed")
	// ErrGetBookByIDNotFound is an error when get not found
	ErrGetBookByIDNotFound = errors.New("get book not found")
	// ErrFailedUpdateBook is an error when update book failed
	ErrFailedUpdateBook = errors.New("update book failed")
	// ErrUpdateBookNotFound is an error when update book not found
	ErrUpdateBookNotFound = errors.New("update book not found")
	// ErrFailedDeleteBook is an error when delete book failed
	ErrFailedDeleteBook = errors.New("delete book failed")
)

// CreateBook creates a new book or updates an existing one based on ISBN
func (l *LibraryService) CreateBook(book *model.CreateBookRequest) error {
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
			"reviewsList",
			"viewsList",
			"wishList"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
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
			"reviewsList" = EXCLUDED."reviewsList",
			"viewsList" = EXCLUDED."viewsList",
			"wishList" = EXCLUDED."wishList",
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
		pq.Array(book.ReviewsList),
		pq.Array(book.ViewsList),
		pq.Array(book.WishList),
	).Scan(&bookID)

	if err != nil {
		log.Error().Msgf("[Error] CreateBook(), db.QueryRow err: %v", err)
		return ErrFailedCreateBook
	}

	return nil
}

// CreateBooksBatch creates multiple books at once
func (l *LibraryService) CreateBooksBatch(books []*model.CreateBookRequest) error {
	if len(books) == 0 {
		return nil // No books to insert
	}

	tx, err := l.db.Begin()
	if err != nil {
		log.Error().Msgf("[Error] CreateBooksBatch(), db.Begin err: %v", err)
		return ErrFailedCreateBook
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
			"reviewsList",
			"viewsList",
			"wishList"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
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
			"reviewsList" = EXCLUDED."reviewsList",
			"viewsList" = EXCLUDED."viewsList",
			"wishList" = EXCLUDED."wishList",
			"updatedAt" = NOW()
		RETURNING "ID";
	`

	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] CreateBooksBatch(), tx.Prepare err: %v", err)
		tx.Rollback()
		return ErrFailedCreateBook
	}
	defer stmt.Close()

	for _, book := range books {
		var bookID string
		err := stmt.QueryRow(
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
			pq.Array(book.ReviewsList),
			pq.Array(book.ViewsList),
			pq.Array(book.WishList),
		).Scan(&bookID)

		if err != nil {
			log.Error().Msgf("[Error] CreateBooksBatch(), stmt.QueryRow err: %v", err)
			tx.Rollback()
			return ErrFailedCreateBook
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Msgf("[Error] CreateBooksBatch(), tx.Commit err: %v", err)
		tx.Rollback()
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
			"updatedAt",
			"reviewsList",
			"viewsList",
			"wishList"
		FROM
			"books"
		WHERE
			"ISBN" = $1;
	`

	var (
		book       model.Book
		updatedAt  sql.NullTime
		reviewList pq.StringArray
		viewList   pq.StringArray
		wishList   pq.StringArray
	)
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
		&updatedAt,
		&reviewList,
		&viewList,
		&wishList,
	)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			log.Error().Msgf("[Error] GetBookByID(), db.QueryRow err: %v", err)
			return nil, ErrGetBookByIDNotFound
		}
		log.Error().Msgf("[Error] GetBookByID(), db.QueryRow err: %v", err)
		return nil, ErrFailedGetBookByID
	}

	book.UpdatedAt = &updatedAt.Time
	book.ReviewsList = reviewList
	book.ViewsList = viewList
	book.WishList = wishList

	return &book, nil
}

// GetBookWithBookID retrieves a book by its ID
func (l *LibraryService) GetBookWithBookID(bookID string) (*model.Book, error) {
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
			"ID" = $1;
	`

	var (
		book       model.Book
		updatedAt  sql.NullTime
		reviewList pq.StringArray
		viewList   pq.StringArray
		wishList   pq.StringArray
	)
	err := l.db.QueryRow(sqlStatement, bookID).Scan(
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
		if errors.Is(sql.ErrNoRows, err) {
			log.Error().Msgf("[Error] GetBookWithBookID(), db.QueryRow err: %v", err)
			return nil, ErrGetBookByIDNotFound
		}
		log.Error().Msgf("[Error] GetBookWithBookID(), db.QueryRow err: %v", err)
		return nil, ErrFailedGetBookByID
	}

	book.UpdatedAt = &updatedAt.Time
	book.ReviewsList = reviewList
	book.ViewsList = viewList
	book.WishList = wishList

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
			log.Error().Msgf("[Error] GetAllBooks(), rows.Scan err: %v", err)
			return nil, err
		}
		book.UpdatedAt = &updatedAt.Time
		book.ReviewsList = reviewList
		book.ViewsList = viewList
		book.WishList = wishList

		books = append(books, book)
	}

	return books, nil
}

// getAllBooksByBookDetailsFrom retrieves all books from the database which match the condition
func (l *LibraryService) GetAllBooksByBookDetailsFrom(request *model.GetAllBooksByBookDetailsFromRequest) ([]model.Book, error) {
	getBookDetailsSqlStatement := `
			SELECT 
				"reservedBookList",
				"pendingBooksList",
				"checkedOutBookList",
				"completedBooksList",
				"wishlistBooks"
			FROM
				"book_details"
			WHERE 
				"userID" = $1`

	var (
		reservedBooksList  pq.StringArray
		pendingBooksList   pq.StringArray
		checkedOutBookList pq.StringArray
		completedBooksList pq.StringArray
		wishlistBooks      pq.StringArray
		bookList           []string
	)
	if err := l.db.QueryRow(getBookDetailsSqlStatement, request.UserID).Scan(&reservedBooksList, &pendingBooksList, &checkedOutBookList, &completedBooksList, &wishlistBooks); err != nil {
		log.Error().Msgf("[Error] GetAllBooks(), db.Query err: %v", err)
		return nil, err
	}

	switch request.BookDetailsFrom {
	case model.BookDetailsFromReserved:
		bookList = reservedBooksList
	case model.BookDetailsFromPending:
		bookList = pendingBooksList
	case model.BookDetailsFromCheckedOut:
		bookList = checkedOutBookList
	case model.BookDetailsFromWishLists:
		bookList = wishlistBooks
	default:
		bookList = wishlistBooks
	}

	bookSqlStatement := `
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
			"ISBN" = ANY($1);
	`

	rows, err := l.db.Query(bookSqlStatement, pq.Array(bookList))
	if err != nil {
		log.Error().Msgf("[Error] GetAllBooks(), db.Query err: %v", err)
		return nil, err
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
			log.Error().Msgf("[Error] GetAllBooks(), rows.Scan err: %v", err)
			return nil, err
		}
		book.UpdatedAt = &updatedAt.Time
		book.ReviewsList = reviewList
		book.ViewsList = viewList
		book.WishList = wishList

		books = append(books, book)
	}

	return books, nil
}

// UpdateBook updates an existing book in the "books" table
func (l *LibraryService) UpdateBook(book *model.UpdateBookRequest) error {
	sqlStatement := `
		UPDATE "books" SET
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
			"updatedAt" = $17,
			"reviewsList" = $18,
			"viewsList" = $19,
			"wishList" = $20
		WHERE
			"ISBN" = $1;
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
		pq.Array(book.ReviewsList),
		pq.Array(book.ViewsList),
		pq.Array(book.WishList),
	)

	if err != nil {
		log.Error().Msgf("[Error] UpdateBook(), db.QueryRow err: %v", err)
		return ErrFailedUpdateBook
	}

	if rowsEffected, err := res.RowsAffected(); err != nil || rowsEffected == 0 {
		log.Error().Msgf("[error] UpdateBook(), [No rows affected]  : %v", err)
		return ErrUpdateBookNotFound
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
