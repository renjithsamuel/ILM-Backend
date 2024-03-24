package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	userIDContextKey = "userID"
	bearerPrefix     = "Bearer "
)

var (
	// ErrTokenExpired is when the token is expired
	ErrTokenExpired = errors.New("token has expired")
)

func (m *UserMiddleware) DoAuthenticate(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.Header("WWW-Authenticate", "Bearer")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized. Bearer token required.",
		})
		c.Abort()
		return
	}

	// Check for Bearer prefix
	if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: token string should start with 'Bearer '",
		})
		c.Abort()
		return
	}

	// Extract token without the Bearer prefix
	token := strings.TrimPrefix(authorizationHeader, bearerPrefix)

	// Validate token using validateToken
	userID, err := m.validateToken(token)
	if err != nil {
		if errors.Is(err, ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		log.Printf("[Validation Failed] %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: " + err.Error(),
		})
		c.Abort()
		return
	}

	if len(userID) == 0 {
		log.Println("[Validation Failed] : UserID not Valid")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "userID not present in token",
		})
		c.Abort()
		return
	}

	c.Set(userIDContextKey, userID)
	c.Next()
}

// validate token
func (m *UserMiddleware) validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKeyByte := []byte(m.secretKey)
		return secretKeyByte, nil
	})
	if err != nil {
		log.Printf("[error] validateToken(): %v\n", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration time
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			log.Println("[error] validateToken(): Token has expired")
			return "", ErrTokenExpired
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			log.Print("[error] validateToken(): sub claim is not a string\n")
			return "", fmt.Errorf("sub claim is not a string")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
