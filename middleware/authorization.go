package middleware

import (
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func SplitToken(headerToken string) string {
	parsToken := strings.SplitAfter(headerToken, " ")
	tokenString := parsToken[1]
	return tokenString
}

func AuthenticateToken(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
}

var ErrUnauthorized = errors.New("Unauthorized")

// AuthenticateToken checks the Authorization token
func AuthorizeToken(c *fiber.Ctx) error {
	headerToken := c.Get("Authorization")

	if headerToken == "" {
		return ErrUnauthorized
	}

	if err := AuthenticateToken(SplitToken(headerToken)); err != nil {
		return ErrUnauthorized
	}
	return nil
}
