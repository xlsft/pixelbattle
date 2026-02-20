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
	"github.com/xlsft/pixelbattle/utils"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(ctx *fiber.Ctx) error {
	db := database.UseDb()
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
	}

	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
	}

	uuidStr, ok := claims["uuid"].(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
	}

	userUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
	}

	var user models.User
	if err := db.Model(&models.User{}).First(&user, "uuid = ?", userUUID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Unauthorized"))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError(err.Error()))
	}

	ctx.Locals("user", user)

	return ctx.Next()
}
