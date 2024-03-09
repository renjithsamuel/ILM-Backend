package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"integrated-library-service/apperror"
	"integrated-library-service/domain"
	"integrated-library-service/model"

	"github.com/rs/zerolog/log"
)

// register user handler creates new user with given data
func (th *LibraryHandler) LoginUserHandler(c *gin.Context) {
	req := model.LoginUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	user, err := th.domain.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrFailedGetUserByEmailNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, err := generateToken(user.UserID, th.secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func generateToken(userID string, secretKey string) (string, error) {
	tokenClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 1000).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	secretKeyByte := []byte(secretKey)
	signedToken, err := token.SignedString(secretKeyByte)
	if err != nil {
		log.Error().Msgf("[error] generateToken: %v", err)
		return "", err
	}

	return signedToken, nil
}
