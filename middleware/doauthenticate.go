package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (m *UserMiddleware) DoAuthenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}

	// Validate token using validateToken
	userID, err := m.validateToken(token)
	if err != nil {
		log.Printf("[Validation Failed] : %v\n", err)
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

	c.Set("userID", userID)
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
		userID, ok := claims["sub"].(string)
		if !ok {
			log.Print("[error] validateToken(): sub claim is not a string\n")
			return "", fmt.Errorf("sub claim is not a string")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
