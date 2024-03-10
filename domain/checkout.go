package domain

import (
	"database/sql"
	"errors"
	"time"

	"integrated-library-service/model"

	"github.com/rs/zerolog/log"
)

var (
	// ErrFailedCreateCheckoutTicket is an error when create checkout failed
	ErrFailedCreateCheckoutTicket = errors.New("create checkout failed")
	// ErrGetCheckoutTicketByIDFailed is an error when get checkout ticket failed
	ErrGetCheckoutTicketByIDFailed = errors.New("get checkout ticket failed")
	// ErrGetCheckoutTicketsFailed is an error when get checkout tickets failed
	ErrGetCheckoutTicketsFailed = errors.New("get checkout tickets failed")
	// ErrGetCheckoutTicketByIDNotFound is an error when get checkout ticket not found
	ErrGetCheckoutTicketByIDNotFound = errors.New("get checkout ticket not found")
	// ErrFailedUpdateCheckoutTicket is an error when update checkout ticket failed
	ErrFailedUpdateCheckoutTicket = errors.New("update checkout ticket failed")
	// ErrFailedDeleteCheckoutTicket is an error when delete checkout ticket not found
	ErrFailedDeleteCheckoutTicket = errors.New("delete checkout ticket failed")
)

// CreateCheckoutTicket creates a new checkout ticket
func (l *LibraryService) CreateCheckoutTicket(ticket *model.CheckoutTicket) error {
	sqlStatement := `
		INSERT INTO "checkout_tickets"(
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
			"createdAt"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING "ID";
	`

	err := l.db.QueryRow(
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
		ticket.CreatedAt,
	).Scan(&ticket.ID)

	if err != nil {
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

	var ticket model.CheckoutTicket
	err := l.db.QueryRow(sqlStatement, ticketID).Scan(
		&ticket.ID,
		&ticket.BookID,
		&ticket.UserID,
		&ticket.IsCheckedOut,
		&ticket.IsReturned,
		&ticket.NumberOfDays,
		&ticket.FineAmount,
		&ticket.ReservedOn,
		&ticket.CheckedOutOn,
		&ticket.ReturnedDate,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msgf("[Error] GetCheckoutTicketByID(), db.QueryRow err: %v", err)
			return nil, ErrGetCheckoutTicketByIDNotFound
		}

		log.Error().Msgf("[Error] GetCheckoutTicketByID(), db.QueryRow err: %v", err)
		return nil, ErrGetCheckoutTicketByIDFailed
	}

	return &ticket, nil
}

// GetAllCheckoutTickets retrieves all checkout tickets
func (l *LibraryService) GetAllCheckoutTickets() ([]model.CheckoutTicket, error) {
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
			"checkout_tickets";
	`

	rows, err := l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetAllCheckoutTickets(), db.Query err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tickets []model.CheckoutTicket
	for rows.Next() {
		var ticket model.CheckoutTicket
		err := rows.Scan(
			&ticket.ID,
			&ticket.BookID,
			&ticket.UserID,
			&ticket.IsCheckedOut,
			&ticket.IsReturned,
			&ticket.NumberOfDays,
			&ticket.FineAmount,
			&ticket.ReservedOn,
			&ticket.CheckedOutOn,
			&ticket.ReturnedDate,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			log.Error().Msgf("[Error] GetAllCheckoutTickets(), rows.Scan err: %v", err)
			return nil, ErrGetCheckoutTicketsFailed
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// UpdateCheckoutTicket updates an existing checkout ticket
func (l *LibraryService) UpdateCheckoutTicket(ticket *model.CheckoutTicket) error {
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