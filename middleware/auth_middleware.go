package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// JwtMiddleware is a Go function that serves as a middleware to handle JWT authentication.
//
// It takes a pointer to a fiber.Ctx as a parameter and returns an error.
func JwtMiddleware(c *fiber.Ctx) error {
	jwtCookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(jwtCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	if claims.ExpiresAt < time.Now().Unix() {
		return c.JSON(fiber.Map{
			"message": "Claims expired",
		})
	}
	return c.Next()
}
