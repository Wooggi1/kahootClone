package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"

	"os"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var secret_key []byte

func init() {
	secret_key = []byte(os.Getenv("JWT_SECRET"))
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		 return secret_key, nil
	})
 
	if err != nil {
		 return err
	}
 
	if !token.Valid {
		 return fmt.Errorf("invalid token")
	}
 
	return nil
}


func JWTAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
	}
	tokenString := parts[1]

	if err := verifyToken(tokenString); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	return c.Next()
}