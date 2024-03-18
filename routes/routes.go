package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	auth "integrated-library-service/middleware"

	"integrated-library-service/handlers"
)

// Route Structure of new routes
type Route struct {
	Name           string
	Method         string
	Pattern        string
	ProtectedRoute bool
	HandlerFunc    gin.HandlerFunc
}

// Routes Array of all available routes
type Routes []Route

// NewRoutes returns the list of available routes
func NewRoutes(libraryHandler *handlers.LibraryHandler) Routes {
	return Routes{
		Route{
			Name:           "Health",
			Method:         http.MethodGet,
			Pattern:        "/health",
			ProtectedRoute: false,
			HandlerFunc:    libraryHandler.HealthHandler,
		},
		// User Related Routes
		Route{
			Name:           "Register New User",
			Method:         http.MethodPost,
			Pattern:        "/users",
			ProtectedRoute: false,
			HandlerFunc:    libraryHandler.RegisterUserHandler,
		},
		Route{
			Name:           "Login User",
			Method:         http.MethodPost,
			Pattern:        "/users/login",
			ProtectedRoute: false,
			HandlerFunc:    libraryHandler.LoginUserHandler,
		},
		Route{
			Name:           "Get User With Book Details",
			Method:         http.MethodGet,
			Pattern:        "/users",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetUserHandler,
		},
		Route{
			Name:           "Get User With Book Details For Given ID",
			Method:         http.MethodGet,
			Pattern:        "/users/:userid",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetUserByIDHandler,
		},
		Route{
			Name:           "Get All Users With Sorted With Book Details",
			Method:         http.MethodGet,
			Pattern:        "/allusers",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetAllUsersHandler,
		},
		Route{
			Name:           "Update User",
			Method:         http.MethodPut,
			Pattern:        "/users",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.UpdateUserHandler,
		},
		Route{
			Name:           "Update User Book Details",
			Method:         http.MethodPut,
			Pattern:        "/users/book-details",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.UpdateBookDetailsHandler,
		},
		Route{
			Name:           "Delete User",
			Method:         http.MethodDelete,
			Pattern:        "/users",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.DeleteUserHandler,
		},
		// book related
		Route{
			Name:           "Create Books in Batch", // will be added manually
			Method:         http.MethodPost,
			Pattern:        "/allbooks",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.CreateBooksBatchHandler,
		},
		Route{
			Name:           "Create Book", // will be added when added to library
			Method:         http.MethodPost,
			Pattern:        "/books",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.CreateBookHandler,
		},
		Route{
			Name:           "Get Book By ISBN",
			Method:         http.MethodGet,
			Pattern:        "/books/:isbn",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetBookByISBNHandler,
		},
		Route{
			Name:           "Get All Books Sorted",
			Method:         http.MethodGet,
			Pattern:        "/allbooks",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetAllBooksHandler,
		},
		Route{
			Name:           "Update Book by ISBN",
			Method:         http.MethodPut,
			Pattern:        "/books",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.UpdateBookHandler,
		},
		Route{
			Name:           "Get All Books By Book Details From",
			Method:         http.MethodGet,
			Pattern:        "/allbooks/:userid/:bookdetailsfrom",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetAllBooksByBookDetailsFromHandler,
		},
		// checkout related
		Route{
			Name:           "Create Checkout Ticket",
			Method:         http.MethodPost,
			Pattern:        "/checkouts",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.CreateCheckoutHandler,
		},
		// Route{
		// 	Name:           "Get Checkout Ticket By ID",
		// 	Method:         http.MethodGet,
		// 	Pattern:        "/checkouts/:checkoutid",
		// 	ProtectedRoute: true,
		// 	HandlerFunc:    libraryHandler.GetCheckoutTicketByIDHandler,
		// },
		Route{
			Name:           "Get Checkout Tickets By BookID and UserID",
			Method:         http.MethodGet,
			Pattern:        "/allcheckouts/:bookid/:userid",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetCheckoutsByUserIDHandler,
		},
		Route{
			Name:           "Get All Checkout Tickets Sorted",
			Method:         http.MethodGet,
			Pattern:        "/allcheckouts",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetAllCheckoutTicketsHandler,
		},
		Route{
			Name:           "Update Checkout Ticket",
			Method:         http.MethodPut,
			Pattern:        "/checkouts",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.UpdateCheckoutTicketHandler,
		},
		Route{
			Name:           "Delete Checkout Ticket",
			Method:         http.MethodDelete,
			Pattern:        "/checkouts/:checkoutid",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.DeleteCheckoutTicketHandler,
		},
		// review related
		Route{
			Name:           "Create Review",
			Method:         http.MethodPost,
			Pattern:        "/reviews",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.CreateReviewHandler,
		},
		Route{
			Name:           "Get Review by ID",
			Method:         http.MethodGet,
			Pattern:        "/reviews/:reviewid",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetReviewByIDHandler,
		},
		Route{
			Name:           "Get Reviews by bookID",
			Method:         http.MethodGet,
			Pattern:        "/allreviews/:bookid",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.GetReviewsByBookIDHandler,
		},
		Route{
			Name:           "Update Review",
			Method:         http.MethodPut,
			Pattern:        "/reviews",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.UpdateReviewHandler,
		},
		Route{
			Name:           "Delete Review",
			Method:         http.MethodDelete,
			Pattern:        "/reviews/:reviewid",
			ProtectedRoute: true,
			HandlerFunc:    libraryHandler.DeleteReviewHandler,
		},
	}
}

// AttachRoutes Attaches routes to the provided server
func AttachRoutes(server *gin.RouterGroup, routes Routes, authMiddleware auth.Middleware) {
	for _, route := range routes {
		if route.ProtectedRoute {
			server.
				Handle(route.Method, route.Pattern, authMiddleware.DoAuthenticate, route.HandlerFunc)
		} else {
			server.
				Handle(route.Method, route.Pattern, route.HandlerFunc)
		}
	}
}
