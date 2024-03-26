package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pyrolass/hotel-reservation-go/entities"
)

func JWTAuthentication(c *fiber.Ctx) error {

	token, ok := c.GetReqHeaders()["Authorization"]

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	tokenClaim, err := verifyToken(strings.Split(token[0], " ")[1])

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims := tokenClaim.(jwt.MapClaims)

	c.Locals("user", claims["id"])

	return c.Next()
}

func GenerateToken(user entities.User) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	expirationTime := time.Now().Add(1 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       user.ID,
			"username": user.FirstName,
			"email":    user.Email,
			"iss":      "hotel-reservation",
			"exp":      expirationTime.Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	secretKey := os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, err
	}

	return token.Claims, nil
}
