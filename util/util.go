package util

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tipbk/sneakfeed-service/model"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSuccessResponse(obj any) map[string]any {
	m := make(map[string]any)
	m["data"] = obj
	return m
}

func GenerateFailedResponse(message any) map[string]any {
	m := make(map[string]any)
	m["error"] = message
	return m
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Generate access token
func GenerateAccessToken(secretString, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"iss":    "SNEAKFEED",
		"exp":    time.Now().Add(time.Minute * 1).Unix(),
	})

	secretKey := []byte(secretString)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate access token
func ValidateAccessToken(secretString, accessToken string) (jwt.MapClaims, error) {
	hmacSecret := []byte(secretString)
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// Validate refresh token
func ValidateRefreshToken(secretString, refreshToken string) (jwt.MapClaims, error) {
	hmacSecret := []byte(secretString)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// Generate refresh token
func GenerateRefreshToken(secretString, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"iss":    "SNEAKFEED",
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	secretKey := []byte(secretString)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserFromContext(c *gin.Context) (*model.User, error) {
	value, ok := c.Get("user")
	if !ok {
		return nil, errors.New("no user")
	}
	return value.(*model.User), nil
}
