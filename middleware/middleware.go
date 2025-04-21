package middleware

import (
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	DoAuthenticate(c *gin.Context)
}

type UserMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *UserMiddleware {
	return &UserMiddleware{
		secretKey: secretKey,
	}
}
