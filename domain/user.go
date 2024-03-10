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
	// ErrFailedCreateUser is an error when create user failed
	ErrFailedCreateUser = errors.New("create user failed")
	// ErrFailedCreateUserBookDetails is an error when create user book details failed
	ErrFailedCreateUserBookDetails = errors.New("create user book details failed")
	// ErrFailedGetUserByEmailFailed is an error when get user failed
	ErrFailedGetUserByEmailFailed = errors.New("get user failed")
	// ErrFailedGetUserByEmailNotFound is an error when get user not found
	ErrFailedGetUserByEmailNotFound = errors.New("get user not found")
	// ErrGetUserWithBookDetailsFailed is an error when get user with book details failed
	ErrGetUserWithBookDetailsFailed = errors.New("get user with book details failed")
	// ErrGetUserWithBookDetailsNotFound is an error when get user with book details NotFound
	ErrGetUserWithBookDetailsNotFound = errors.New("get user with book details not found")
	// ErrFailedUpdateUser is an error when update user failed
	ErrFailedUpdateUser = errors.New("update user failed")
	// ErrFailedUpdateBookDetails is an error when update book details failed
	ErrFailedUpdateBookDetails = errors.New("update book details failed")
	// ErrFailedDeleteUser is an error when delete user failed
	ErrFailedDeleteUser = errors.New("delete user failed")
)

// create user creates new user
func (l *LibraryService) CreateUser(user *model.RegisterUserRequest) error {
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
								"role" = EXCLUDED."role"
					RETURNING "userID";
					`
	var userID string
	if err := l.db.QueryRow(sqlStatement, user.ProfileImageUrl, user.Name, user.Email, user.Role).Scan(&userID); err != nil {
		log.Error().Msgf("[Error] CreateUser(), db.Exec err: %v", err)
		return ErrFailedCreateUser
	}

	if err := l.createUserBookDetails(userID); err != nil {
		log.Error().Msgf("[Error] CreateUser(), db.Exec err: %v", err)
		return err
	}

	return nil
}

// CreateUserBookDetails user creates new user book details
func (l *LibraryService) createUserBookDetails(userID string) error {
	sqlStatement := `
						INSERT INTO "book_details"(
									"userID"
								) VALUES ($1);
						;
					`

	if _, err := l.db.Exec(sqlStatement, userID); err != nil {
		log.Error().Msgf("[Error] createUserBookDetails(), db.Exec err: %v", err)
		return ErrFailedCreateUserBookDetails
	}

	return nil
}

// GetUserByEmail gets user with unique emailID
func (l *LibraryService) GetUserByEmail(email string) (*model.User, error) {
	sqlStatement := `
		SELECT 
			"userID",
			"profileImageUrl",
			"name",
			"email",
			"role",
			"dateOfBirth",
			"phoneNumber",
			"address",
			"joinedDate",
			"country",
			"views",
			"fineAmount",
			"isPaymentDone",
			"createdAt",
			"updatedAt"   
		FROM 
			"users"
		WHERE 
			"email" = $1;
	`

	var user model.User
	err := l.db.QueryRow(sqlStatement, email).Scan(
		&user.UserID,
		&user.ProfileImageUrl,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.DateOfBirth,
		&user.PhoneNumber,
		&user.Address,
		&user.JoinedDate,
		&user.Country,
		&user.Views,
		&user.FineAmount,
		&user.IsPaymentDone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msgf("[Error] GetUserByEmail(), db.QueryRow err: %v", err)
			return nil, ErrFailedGetUserByEmailNotFound
		}

		log.Error().Msgf("[Error] GetUserByEmail(), db.QueryRow err: %v", err)
		return nil, ErrFailedGetUserByEmailFailed
	}

	return &user, nil
}

// GetUserWithBookDetails gets user with userID and gives user data with book details
func (l *LibraryService) GetUserWithBookDetails(userID string) (*model.User, error) {
	sqlStatement := `
		SELECT 
			u."userID",
			u."profileImageUrl",
			u."name",
			u."email",
			u."role",
			u."dateOfBirth",
			u."phoneNumber",
			u."address",
			u."joinedDate",
			u."country",
			u."views",
			u."fineAmount",
			u."isPaymentDone",
			u."createdAt",
			u."updatedAt",
			bkd."reservedBooksCount",
			bkd."reservedBookList",
			bkd."pendingBooksCount",
			bkd."pendingBooksList",
			bkd."checkedOutBooksCount",
			bkd."checkedOutBookList",
			bkd."completedBooksCount",
			bkd."completedBooksList",
			bkd."favoriteGenres",
			bkd."wishlistBooks",
			bkd."createdAt" as "bookDetails.createdAt",
			bkd."updatedAt" as "bookDetails.updatedAt"
		FROM 
			"users" as u INNER JOIN "book_details" as bkd 
				ON u."userID" = bkd."userID"
		WHERE 
			u."userID" = $1 AND
			bkd."userID" = $1
	`

	var (
		user               model.User
		reservedBooksList  pq.StringArray
		pendingBooksList   pq.StringArray
		checkedOutBookList pq.StringArray
		completedBooksList pq.StringArray
		favoriteGenres     pq.StringArray
		wishlistBooks      pq.StringArray
	)
	err := l.db.QueryRow(sqlStatement, userID).Scan(
		&user.UserID,
		&user.ProfileImageUrl,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.DateOfBirth,
		&user.PhoneNumber,
		&user.Address,
		&user.JoinedDate,
		&user.Country,
		&user.Views,
		&user.FineAmount,
		&user.IsPaymentDone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.BookDetails.ReservedBooksCount,
		&reservedBooksList,
		&user.BookDetails.PendingBooksCount,
		&pendingBooksList,
		&user.BookDetails.CheckedOutBooksCount,
		&checkedOutBookList,
		&user.BookDetails.CompletedBooksCount,
		&completedBooksList,
		&favoriteGenres,
		&wishlistBooks,
		&user.BookDetails.CreatedAt,
		&user.BookDetails.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msgf("[Error] GetUserByEmail(), db.QueryRow err: %v", err)
			return nil, ErrFailedGetUserByEmailNotFound
		}

		log.Error().Msgf("[Error] GetUserByEmail(), db.QueryRow err: %v", err)
		return nil, ErrFailedGetUserByEmailFailed
	}

	user.BookDetails.ReservedBookList = reservedBooksList
	user.BookDetails.PendingBooksList = pendingBooksList
	user.BookDetails.CheckedOutBookList = checkedOutBookList
	user.BookDetails.CompletedBooksList = completedBooksList
	user.BookDetails.WishlistBooks = wishlistBooks
	// Convert pq.StringArray to []model.BookGenreType
	user.BookDetails.FavoriteGenres = make([]model.BookGenreType, len(favoriteGenres))
	for i, genreString := range favoriteGenres {
		user.BookDetails.FavoriteGenres[i] = model.BookGenreType("").StringToBookGenre(genreString)
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (l *LibraryService) GetAllUsers() ([]model.User, error) {
	sqlStatement := `
		SELECT 
			"userID",
			"profileImageUrl",
			"name",
			"email",
			"role",
			"dateOfBirth",
			"phoneNumber",
			"address",
			"joinedDate",
			"country",
			"views",
			"fineAmount",
			"isPaymentDone",
			"createdAt",
			"updatedAt"
		FROM 
			"users";
	`

	rows, err := l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetAllUsers(), db.Query err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.UserID,
			&user.ProfileImageUrl,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.DateOfBirth,
			&user.PhoneNumber,
			&user.Address,
			&user.JoinedDate,
			&user.Country,
			&user.Views,
			&user.FineAmount,
			&user.IsPaymentDone,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Error().Msgf("[Error] GetAllUsers(), rows.Scan err: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}


// UpdateUser updates an existing user in the "users" table
func (l *LibraryService) UpdateUser(user *model.User, userID string) error {
	sqlStatement := `
		UPDATE "users" SET
			"email" = $1,
			"profileImageUrl" = $2,
			"name" = $3,
			"role" = $4,
			"dateOfBirth" = $5,
			"phoneNumber" = $6,
			"address" = $7,
			"joinedDate" = $8,
			"country" = $9,
			"views" = $10,
			"fineAmount" = $11,
			"isPaymentDone" = $12,
			"updatedAt" = $13
		WHERE
			"userID" = $14;
	`
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	res, err := l.db.Exec(
		sqlStatement,
		user.Email,
		user.ProfileImageUrl,
		user.Name,
		user.Role,
		user.DateOfBirth,
		user.PhoneNumber,
		user.Address,
		user.JoinedDate,
		user.Country,
		user.Views,
		user.FineAmount,
		user.IsPaymentDone,
		updatedAt,
		userID,
	)

	if err != nil {
		log.Error().Msgf("[Error] UpdateUser(), db.QueryRow err: %v", err)
		return ErrFailedUpdateUser
	}

	if rowsEffected, err := res.RowsAffected(); err != nil || rowsEffected == 0 {
		log.Error().Msgf("[error] UpdateUser(), [No rows affected]  : %v", err)
		return ErrFailedUpdateUser
	}

	return nil
}

// UpdateBookDetails updates the "book_details" table for a specific user
func (l *LibraryService) UpdateBookDetails(bookDetails *model.BookDetails, userID string) error {
	sqlStatement := `
		UPDATE "book_details" SET
			"reservedBooksCount" = $2,
			"reservedBookList" = $3,
			"checkedOutBooksCount" = $4,
			"checkedOutBookList" = $5,
			"pendingBooksCount" = $6,
			"pendingBooksList" = $7,
			"completedBooksCount" = $8,
			"completedBooksList" = $9,
			"favoriteGenres" = $10,
			"wishlistBooks" = $11,
			"updatedAt" = $12
		WHERE
			"userID" = $1;
	`

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	res, err := l.db.Exec(
		sqlStatement,
		userID,
		bookDetails.ReservedBooksCount,
		pq.Array(bookDetails.ReservedBookList),
		bookDetails.CheckedOutBooksCount,
		pq.Array(bookDetails.CheckedOutBookList),
		bookDetails.PendingBooksCount,
		pq.Array(bookDetails.PendingBooksList),
		bookDetails.CompletedBooksCount,
		pq.Array(bookDetails.CompletedBooksList),
		pq.Array(bookDetails.FavoriteGenres),
		pq.Array(bookDetails.WishlistBooks),
		updatedAt,
	)
	if err != nil {
		log.Error().Msgf("[Error] UpdateBookDetails(), db.QueryRow err: %v", err)
		return ErrFailedUpdateBookDetails
	}

	if rowsEffected, err := res.RowsAffected(); err != nil || rowsEffected == 0 {
		log.Error().Msgf("[Error] UpdateBookDetails(), [No rows affected]  : %v", err)
		return ErrFailedUpdateBookDetails
	}

	return nil
}

// DeleteUser deletes a user from the "users" table based on userID
func (l *LibraryService) DeleteUser(userID string) error {
	sqlStatement := `
		DELETE FROM "users" WHERE "userID" = $1;
	`

	_, err := l.db.Exec(sqlStatement, userID)
	if err != nil {
		log.Error().Msgf("[Error] DeleteUser(), db.Exec err: %v", err)
		return ErrFailedDeleteUser
	}

	return nil
}
