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
			Name:           "Get All Users With Sorted",
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
			Name:           "Create Book", // will be added when added to library
			Method:         http.MethodGet,
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
		// checkout related
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
