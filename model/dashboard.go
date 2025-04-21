package model

import "time"

// DashboardLineGraphData
type DashboardLineGraphData struct {
	// month
	Month time.Time `json:"month" binding:"required"`
	// NoOfActiveusers is the number of Active users of the application with membership in this month
	NoOfActiveusers int `json:"noOfActiveUsers" binding:"required"`
	// NoOfCheckouts is the number of Checkouts in this month
	NoOfCheckouts int `json:"noOfCheckouts" binding:"required"`
}

// DashboardDataBoard
type DashboardDataBoard struct {
	// each of the below counts is for total not per month/year - its from the whole DB
	UsersCount     int `json:"usersCount" binding:"required"`
	BooksCount     int `json:"booksCount" binding:"required"`
	CheckoutsCount int `json:"checkoutsCount" binding:"required"`
	// Revenue is collection of Total fineAmount and for each user with isPaymentDone * 30 in rupees
	RevenueAmountTotal int `json:"revenueAmount" binding:"required"`
	// monthly data
	// MonthlyNewBooksAddedCount nothing but book with createdAt within this month
	MonthlyNewBooksAddedCount int `json:"monthlyNewBooksAddedCount" binding:"required"`
	// MonthlyNewRegisteredUserCount nothing but user with createdAt within this month
	MonthlyNewRegisteredUserCount int `json:"monthlyNewRegisteredUserCount" binding:"required"`
	//MonthlyNewCheckoutTicketsCount nothing but checkout_tickets with createdAt within this month
	MonthlyNewCheckoutTicketsCount int `json:"monthlyNewCheckoutTicketsCount" binding:"required"`
	// MonthlyFineAmountTotal nothing but sum of fineAmount which is in the checkout_tickets with createdAt within this month
	MonthlyFineAmountTotal int `json:"monthlyFineAmountTotal" binding:"required"`
}

// HighDemandBooks is books sorted based on wishListCount and Limit 3
type HighDemandBooks []Book
