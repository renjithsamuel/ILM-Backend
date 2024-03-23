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
	// ErrGetDashboardLineGraphDataFailed is when get dashboard graph data fails
	ErrGetDashboardLineGraphDataFailed = errors.New("get dashboard graph data failed")

	// ErrGetDashboardDataBoardFailed is when get dashboard data board fails
	ErrGetDashboardDataBoardFailed = errors.New("get dashboard data board failed")

	// ErrGetHighDemandBooksFailed is when get high demand book fails
	ErrGetHighDemandBooksFailed = errors.New("get high demand book failed")
)

// GetDashboardLineGraphData retrieves Dashboard's Graph related data
func (l *LibraryService) GetDashboardLineGraphData() ([]model.DashboardLineGraphData, error) {
	// dashboardData is a collection of past 7 Months, ordered on date in ascending order
	dashboardData := make([]model.DashboardLineGraphData, 0)

	// Query to generate a series of the past 7 months
	seriesSQL := `
		SELECT 
			DATE_TRUNC('month', NOW() - INTERVAL '6 months' + INTERVAL '1 month' * generate_series(0, 6)) AS month
	`

	// Query to retrieve dashboardData for the past 7 months
	sqlStatement := `
		SELECT 
			series.month,
			COALESCE(COUNT(DISTINCT u."userID"), 0) AS noOfActiveUsers,
			COALESCE(COUNT(c."ID"), 0) AS noOfCheckouts
		FROM 
			(SELECT 
				DATE_TRUNC('month', NOW() - INTERVAL '6 months' + INTERVAL '1 month' * generate_series(0, 6)) AS month
			) AS series
		LEFT JOIN 
			"users" AS u ON DATE_TRUNC('month', u."createdAt") = series.month
		LEFT JOIN 
			"checkout_tickets" AS c ON DATE_TRUNC('month', c."createdAt") = series.month AND u."userID" = c."userID"
		GROUP BY 
			1
		ORDER BY 
			1;
	`

	// Execute the series query to get all months
	rows, err := l.db.Query(seriesSQL)
	if err != nil {
		log.Error().Msgf("[Error] GetDashboardLineGraphData(), series query err: %v", err)
		return nil, ErrGetDashboardLineGraphDataFailed
	}
	defer rows.Close()

	// Iterate over the series of months
	for rows.Next() {
		var month time.Time
		err := rows.Scan(&month)
		if err != nil {
			log.Error().Msgf("[Error] GetDashboardLineGraphData(), series rows.Scan err: %v", err)
			return nil, ErrGetDashboardLineGraphDataFailed
		}

		// Append the month with no data to the dashboardData
		dashboardData = append(dashboardData, model.DashboardLineGraphData{
			Month:           month,
			NoOfActiveusers: 0,
			NoOfCheckouts:   0,
		})
	}

	// Execute the main query to retrieve data for each month
	rows, err = l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetDashboardLineGraphData(), db.Query err: %v", err)
		return nil, ErrGetDashboardLineGraphDataFailed
	}
	defer rows.Close()

	for rows.Next() {
		var dashboardLineGraphData model.DashboardLineGraphData
		err := rows.Scan(
			&dashboardLineGraphData.Month,
			&dashboardLineGraphData.NoOfActiveusers,
			&dashboardLineGraphData.NoOfCheckouts,
		)
		if err != nil {
			log.Error().Msgf("[Error] GetDashboardLineGraphData(), rows.Scan err: %v", err)
			return nil, ErrGetDashboardLineGraphDataFailed
		}

		// Update the corresponding month's data
		for i := range dashboardData {
			if dashboardData[i].Month.Equal(dashboardLineGraphData.Month) {
				dashboardData[i].NoOfActiveusers = dashboardLineGraphData.NoOfActiveusers
				dashboardData[i].NoOfCheckouts = dashboardLineGraphData.NoOfCheckouts
				break
			}
		}
	}

	return dashboardData, nil
}

// GetDashboardDataBoard retrieves for Dashboard's data board
func (l *LibraryService) GetDashboardDataBoard() (*model.DashboardDataBoard, error) {
	var dashboardData model.DashboardDataBoard

	// Query to get total counts from the entire database
	totalCountsSQL := `
		SELECT
			(SELECT COUNT(*) FROM users) AS usersCount,
			(SELECT COUNT(*) FROM books) AS booksCount,
			(SELECT COUNT(*) FROM checkout_tickets) AS checkoutsCount,
			(SELECT SUM("fineAmount") FROM checkout_tickets) AS revenueAmount
	`

	// Query to get monthly counts for the current month
	monthlyCountsSQL := `
		SELECT
			(SELECT COUNT(*) FROM books WHERE DATE_TRUNC('month', "createdAt") = DATE_TRUNC('month', CURRENT_DATE)) AS monthlyNewBooksAddedCount,
			(SELECT COUNT(*) FROM users WHERE DATE_TRUNC('month', "createdAt") = DATE_TRUNC('month', CURRENT_DATE)) AS monthlyNewRegisteredUserCount,
			(SELECT COUNT(*) FROM checkout_tickets WHERE DATE_TRUNC('month', "createdAt") = DATE_TRUNC('month', CURRENT_DATE)) AS monthlyNewCheckoutTicketsCount,
			(SELECT SUM("fineAmount") FROM checkout_tickets WHERE DATE_TRUNC('month', "createdAt") = DATE_TRUNC('month', CURRENT_DATE)) AS monthlyFineAmountTotal
	`

	// Retrieve total counts
	err := l.db.QueryRow(totalCountsSQL).Scan(
		&dashboardData.UsersCount,
		&dashboardData.BooksCount,
		&dashboardData.CheckoutsCount,
		&dashboardData.RevenueAmountTotal,
	)
	if err != nil {
		log.Error().Msgf("[Error] GetDashboardDataBoard(), retrieving total counts: %v", err)
		return &model.DashboardDataBoard{}, ErrGetDashboardDataBoardFailed
	}

	// Retrieve monthly counts for the current month
	err = l.db.QueryRow(monthlyCountsSQL).Scan(
		&dashboardData.MonthlyNewBooksAddedCount,
		&dashboardData.MonthlyNewRegisteredUserCount,
		&dashboardData.MonthlyNewCheckoutTicketsCount,
		&dashboardData.MonthlyFineAmountTotal,
	)
	if err != nil {
		log.Error().Msgf("[Error] GetDashboardDataBoard(), retrieving monthly counts: %v", err)
		return &model.DashboardDataBoard{}, ErrGetDashboardDataBoardFailed
	}

	return &dashboardData, nil
}

// GetHighDemandBooks retrieves high demand books sorted by wishListCount and limited to 3
func (l *LibraryService) GetHighDemandBooks() (*model.HighDemandBooks, error) {
	var highDemandBooks model.HighDemandBooks

	// Query to retrieve high demand books
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
        ORDER BY 
            "wishlistCount" DESC
        LIMIT 3;
    `

	rows, err := l.db.Query(sqlStatement)
	if err != nil {
		log.Error().Msgf("[Error] GetHighDemandBooks(), db.Query err: %v", err)
		return nil, ErrGetHighDemandBooksFailed
	}
	defer rows.Close()

	for rows.Next() {
		var book model.Book
		var updatedAt sql.NullTime
		var reviewList, viewList, wishList pq.StringArray

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
			log.Error().Msgf("[Error] GetHighDemandBooks(), rows.Scan err: %v", err)
			return nil, ErrGetHighDemandBooksFailed
		}
		book.UpdatedAt = &updatedAt.Time
		book.ReviewsList = reviewList
		book.ViewsList = viewList
		book.WishList = wishList

		highDemandBooks = append(highDemandBooks, book)
	}

	return &highDemandBooks, nil
}
