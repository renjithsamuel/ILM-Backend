package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"integrated-library-service/model"

	"github.com/rs/zerolog/log"
)

var (
	// ErrFailedCreateCheckoutTicket is an error when create checkout failed
	ErrFailedCreateCheckoutTicket = errors.New("create checkout failed")
	// ErrFailedCreateCheckoutTicketConflict is an error when create checkout failed conflict
	ErrFailedCreateCheckoutTicketConflict = errors.New("create checkout failed with conflict")
	// ErrGetCheckoutTicketByIDFailed is an error when get checkout ticket failed
	ErrGetCheckoutTicketByIDFailed = errors.New("get checkout ticket failed")
	// ErrGetCheckoutTicketsFailed is an error when get checkout tickets failed
	ErrGetCheckoutTicketsFailed = errors.New("get checkout tickets failed")
	// ErrGetCheckoutTicketByIDNotFound is an error when get checkout ticket not found
	ErrGetCheckoutTicketByIDNotFound = errors.New("get checkout ticket not found")
	// ErrGetCheckoutByUserIDNotFailed is an error when get checkout by userid ticket failed
	ErrGetCheckoutByUserIDNotFailed = errors.New("get checkout by userid ticket failed")
	// ErrGetCheckoutByUserIDNotFound is an error when get checkout by userid ticket not found
	ErrGetCheckoutByUserIDNotFound = errors.New("get checkout by userid not found")
	// ErrFailedUpdateCheckoutTicket is an error when update checkout ticket failed
	ErrFailedUpdateCheckoutTicket = errors.New("update checkout ticket failed")
	// ErrFailedDeleteCheckoutTicket is an error when delete checkout ticket not found
	ErrFailedDeleteCheckoutTicket = errors.New("delete checkout ticket failed")
	// ErrPaymentPending is an error when payment is pending
	ErrPaymentPending = errors.New("payment is pending")
	// ErrOutOfStock is an error when book is out of stock
	ErrOutOfStock = errors.New("book is out of stock")
)

// CreateCheckoutTicket creates a new checkout ticket
func (l *LibraryService) CreateCheckoutTicket(ticket *model.CreateCheckoutRequest) error {
	user, err := l.GetUserWithBookDetails(ticket.UserID)
	if err != nil {
		log.Error().Msgf("[Error] CreateCheckoutTicket(), GetUserWithBookDetails err: %v", err)
		return ErrFailedCreateCheckoutTicket
	}

	if !user.IsPaymentDone {
		return ErrPaymentPending
	}

	book, err := l.GetBookWithBookID(ticket.BookID)
	if err != nil {
		log.Error().Msgf("[Error] CreateCheckoutTicket(), GetBookWithBookID err: %v", err)
		return ErrFailedCreateCheckoutTicket
	}

	if book.BooksLeft == 0 {
		return ErrOutOfStock
	}

	// if other checkouts which has the same book ID and userID is there, he is spamming
	otherCheckouts, err := l.GetCheckoutsByUserID(ticket.BookID, ticket.UserID)
	if err != nil {
		log.Error().Msgf("[Error] CreateCheckoutTicket(), GetCheckoutsByUserID err: %v", err)
		return ErrFailedCreateCheckoutTicket
	}

	for _, checkout := range otherCheckouts {
		if !checkout.IsReturned {
			log.Error().Msgf("[Error] CreateCheckoutTicket(), Book is already Pending Return")
			return ErrFailedCreateCheckoutTicketConflict
		}
	}

	sqlStatement := `
		INSERT INTO "checkout_tickets"(
			"bookID",
			"userID",
			"numberOfDays",
			"reservedOn"
		) VALUES (
			$1, $2, $3, $4
		);
	`

	updateAt := time.Now().UTC().Format(time.RFC1123)
	res := l.db.QueryRow(
		sqlStatement,
		ticket.BookID,
		ticket.UserID,
		ticket.NumberOfDays,
		updateAt,
	)

	if err := res.Err(); err != nil {
		log.Error().Msgf("[Error] CreateCheckoutTicket(), db.QueryRow err: %v", err)
		return ErrFailedCreateCheckoutTicket
	}

	return nil
}

// GetCheckoutTicketByID retrieves a checkout ticket by its ID
func (l *LibraryService) GetCheckoutTicketByID(ticketID string) (*model.CheckoutTicket, error) {
	sqlStatement := `
		SELECT 
			"ID",
			"bookID",
			"userID",
			"isCheckedOut",
			"isReturned",
			"numberOfDays",
			"fineAmount",
			"reservedOn",
			"checkedOutOn",
			"returnedDate",
			"createdAt",
			"updatedAt"
		FROM 
			"checkout_tickets"
		WHERE 
			"ID" = $1;
	`

	var (
		ticket       model.CheckoutTicket
		updatedAt    sql.NullTime
		reservedOn   sql.NullTime
		checkedOutOn sql.NullTime
		returnedDate sql.NullTime
	)
	err := l.db.QueryRow(sqlStatement, ticketID).Scan(
		&ticket.ID,
		&ticket.BookID,
		&ticket.UserID,
		&ticket.IsCheckedOut,
		&ticket.IsReturned,
		&ticket.NumberOfDays,
		&ticket.FineAmount,
		&reservedOn,
		&checkedOutOn,
		&returnedDate,
		&ticket.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msgf("[Error] GetCheckoutTicketByID(), db.QueryRow err: %v", err)
			return nil, ErrGetCheckoutTicketByIDNotFound
		}

		log.Error().Msgf("[Error] GetCheckoutTicketByID(), db.QueryRow err: %v", err)
		return nil, ErrGetCheckoutTicketByIDFailed
	}
	ticket.UpdatedAt = updatedAt.Time
	ticket.CheckedOutOn = checkedOutOn.Time
	ticket.ReturnedDate = returnedDate.Time
	ticket.ReservedOn = reservedOn.Time

	return &ticket, nil
}

// GetCheckoutsByUserID retrieves a checkout tickets by its UserID and BookID
func (l *LibraryService) GetCheckoutsByUserID(bookID, userID string) ([]model.CheckoutTicket, error) {
	sqlStatement := `
		SELECT 
			"ID",
			"bookID",
			"userID",
			"isCheckedOut",
			"isReturned",
			"numberOfDays",
			"fineAmount",
			"reservedOn",
			"checkedOutOn",
			"returnedDate",
			"createdAt",
			"updatedAt"
		FROM 
			"checkout_tickets"
		WHERE 
			"bookID" = $1 AND "userID" = $2;
	`

	rows, err := l.db.Query(sqlStatement, bookID, userID)
	if err != nil {
		log.Error().Msgf("[Error] GetCheckoutsTicketByID(), db.QueryRow err: %v", err)
		return nil, ErrGetCheckoutTicketByIDFailed
	}
	var tickets []model.CheckoutTicket
	for rows.Next() {
		var (
			ticket       model.CheckoutTicket
			updatedAt    sql.NullTime
			reservedOn   sql.NullTime
			checkedOutOn sql.NullTime
			returnedDate sql.NullTime
		)
		err := rows.Scan(
			&ticket.ID,
			&ticket.BookID,
			&ticket.UserID,
			&ticket.IsCheckedOut,
			&ticket.IsReturned,
			&ticket.NumberOfDays,
			&ticket.FineAmount,
			&reservedOn,
			&checkedOutOn,
			&returnedDate,
			&ticket.CreatedAt,
			&updatedAt,
		)

		if err != nil {
			log.Error().Msgf("[Error] GetCheckoutsByUserID(), rows.Scan err: %v", err)
			return nil, ErrGetCheckoutTicketByIDFailed
		}
		ticket.UpdatedAt = updatedAt.Time
		ticket.CheckedOutOn = checkedOutOn.Time
		ticket.ReturnedDate = returnedDate.Time
		ticket.ReservedOn = reservedOn.Time
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// GetAllCheckoutTicketsWithDetails retrieves all checkout tickets with associated user and book data
func (l *LibraryService) GetAllCheckoutTicketsWithDetails(request *model.GetAllCheckoutData) ([]model.CheckoutTicketResponse, uint, error) {
	sqlStatement := `
		SELECT 
			ct."ID",
			ct."bookID",
			ct."userID",
			ct."isCheckedOut",
			ct."isReturned",
			ct."numberOfDays",
			ct."fineAmount",
			ct."reservedOn",
			ct."checkedOutOn",
			ct."returnedDate",
			ct."createdAt",
			ct."updatedAt",
			u."profileImageUrl" as "userProfileImageUrl",
			u."name" as "userName",
			u."email" as "userEmail",
			u."role" as "userRole",
			u."dateOfBirth" as "userDateOfBirth",
			u."phoneNumber" as "userPhoneNumber",
			u."address" as "userAddress",
			u."joinedDate" as "userJoinedDate",
			u."country" as "userCountry",
			u."views" as "userViews",
			u."fineAmount" as "userFineAmount",
			b."ID" as "bookID",
			b."ISBN" as "bookISBN",
			b."title" as "bookTitle",
			b."author" as "bookAuthor",
			b."genre" as "bookGenre",
			b."publishedDate" as "bookPublishedDate",
			b."desc" as "bookDescription",
			b."previewLink" as "bookPreviewLink",
			b."coverImage" as "bookCoverImage",
			b."shelfNumber" as "bookShelfNumber",
			b."inLibrary" as "bookInLibrary",
			b."views" as "bookViews",
			b."booksLeft" as "bookBooksLeft",
			b."wishlistCount" as "bookWishlistCount",
			b."rating" as "bookRating",
			b."reviewCount" as "bookReviewCount",
			b."approximateDemand" as "bookApproximateDemand",
			b."createdAt" as "bookCreatedAt",
			b."updatedAt" as "bookUpdatedAt"
		FROM 
			"checkout_tickets" ct
		INNER JOIN
			"users" u ON ct."userID" = u."userID"
		INNER JOIN
			"books" b ON ct."bookID" = b."ID"
		ORDER BY 
			%s -- orderby
		%s; -- criteria for limit and offset 	
		;
	`

	orderBy := `%s ASC`

	if request.OrderBy == "descending" {
		orderBy = `%s DESC`
	}

	switch request.SortBy {
	case "reservedOn":
		orderBy = fmt.Sprintf(orderBy, `ct."reservedOn"`)
	case "checkedoutOn":
		orderBy = fmt.Sprintf(orderBy, `ct."checkedOutOn"`)
	case "returnedOn":
		orderBy = fmt.Sprintf(orderBy, `ct."returnedDate"`)
	case "fineAmount":
		orderBy = fmt.Sprintf(orderBy, `ct."fineAmount"`)
	default:
		orderBy = fmt.Sprintf(orderBy, `ct."reservedOn"`)
	}

	limitOffset := ` LIMIT %d OFFSET %d`
	limitOffset = fmt.Sprintf(limitOffset, request.Limit, (request.Page-1)*(request.Limit))
	sqlStatement = fmt.Sprintf(sqlStatement, orderBy, limitOffset)

	rows, err := l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetAllCheckoutTicketsWithDetails(), db.Query err: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var tickets []model.CheckoutTicketResponse
	for rows.Next() {
		var ticket model.CheckoutTicketResponse
		var user model.User
		var book model.Book
		var checkedOutOn sql.NullTime
		var returnedDate sql.NullTime
		var updatedAt sql.NullTime
		var reservedOn sql.NullTime

		err := rows.Scan(
			&ticket.ID,
			&ticket.BookID,
			&ticket.UserID,
			&ticket.IsCheckedOut,
			&ticket.IsReturned,
			&ticket.NumberOfDays,
			&ticket.FineAmount,
			&reservedOn,
			&checkedOutOn,
			&returnedDate,
			&ticket.CreatedAt,
			&updatedAt,
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
			log.Error().Msgf("[Error] GetAllCheckoutTicketsWithDetails(), rows.Scan err: %v", err)
			return nil, 0, ErrGetCheckoutTicketsFailed
		}

		ticket.User = user
		ticket.Book = book
		ticket.CheckedOutOn = checkedOutOn.Time
		ticket.ReturnedDate = returnedDate.Time
		ticket.UpdatedAt = updatedAt.Time
		ticket.ReservedOn = reservedOn.Time

		tickets = append(tickets, ticket)
	}

	sqlStatementCount := `
	SELECT 
		COUNT(*)
	FROM 
		"checkout_tickets" ct
	INNER JOIN
		"users" u ON ct."userID" = u."userID"
	INNER JOIN
		"books" b ON ct."bookID" = b."ID"
	%s; --LIMIT AND OFFSET 
`

	var totalRows uint
	err = l.db.QueryRow(fmt.Sprintf(sqlStatementCount, limitOffset)).Scan(&totalRows)
	// no rows
	if errors.Is(err, sql.ErrNoRows) {
		return []model.CheckoutTicketResponse{}, 0, nil
	}
	if err != nil {
		log.Error().Msgf("[Error] GetAllCheckoutTicketsWithDetails(), count query err: %v", err)
		return nil, 0, err
	}
	// Calculate total pages
	totalPages := (uint32(totalRows) + request.Limit - 1) / request.Limit

	return tickets, uint(totalPages), nil
}

// UpdateCheckoutTicket updates an existing checkout ticket
func (l *LibraryService) UpdateCheckoutTicket(ticket *model.UpdateCheckoutTicketRequest) error {
	sqlStatement := `
		UPDATE "checkout_tickets" SET
			"bookID" = $2,
			"userID" = $3,
			"isCheckedOut" = $4,
			"isReturned" = $5,
			"numberOfDays" = $6,
			"fineAmount" = $7,
			"reservedOn" = $8,
			"checkedOutOn" = $9,
			"returnedDate" = $10,
			"updatedAt" = $11
		WHERE
			"ID" = $1;
	`

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	_, err := l.db.Exec(
		sqlStatement,
		ticket.ID,
		ticket.BookID,
		ticket.UserID,
		ticket.IsCheckedOut,
		ticket.IsReturned,
		ticket.NumberOfDays,
		ticket.FineAmount,
		ticket.ReservedOn,
		ticket.CheckedOutOn,
		ticket.ReturnedDate,
		updatedAt,
	)

	if err != nil {
		log.Error().Msgf("[Error] UpdateCheckoutTicket(), db.Exec err: %v", err)
		return ErrFailedUpdateCheckoutTicket
	}

	return nil
}

// DeleteCheckoutTicket deletes a checkout ticket by its ID
func (l *LibraryService) DeleteCheckoutTicket(ticketID string) error {
	sqlStatement := `
		DELETE FROM "checkout_tickets" WHERE "ID" = $1;
	`

	_, err := l.db.Exec(sqlStatement, ticketID)
	if err != nil {
		log.Error().Msgf("[Error] DeleteCheckoutTicket(), db.Exec err: %v", err)
		return ErrFailedDeleteCheckoutTicket
	}

	return nil
}
