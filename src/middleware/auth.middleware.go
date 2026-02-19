package middleware

import (
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() fiber.Handler {
	db := database.UseDb()
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		uuidStr, ok := claims["uuid"].(string)
		if !ok {
			return fiber.ErrUnauthorized
		}

		userUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		var user models.UserModel
		if err := db.Model(&models.UserModel{}).First(&user, "uuid = ?", userUUID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fiber.ErrUnauthorized
			}
			return fiber.ErrInternalServerError
		}

		c.Locals("user", user)

		return c.Next()
	}
}
