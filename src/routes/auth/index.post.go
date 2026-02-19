package authRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
	"github.com/xlsft/pixelbattle/utils"
)

func HandlePost(ctx *fiber.Ctx) error {
	request := utils.TelegramData{}
	db := database.UseDb()

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	if err := request.VerifyTelegramData(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	var user models.User

	if err := db.Model(&models.User{}).Where("id = ?", request.ID).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Database error"))
		}

		user = models.User{
			UUID:     uuid.New(),
			ID:       request.ID,
			Name:     request.FirstName + " " + request.LastName,
			Nickname: request.Username,
			Picture:  request.PhotoURL,
		}

		if err := db.Model(&models.User{}).Create(&user).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Failed to create user"))
		}
	}

	token, err := utils.GenerateJWT(user.UUID.String())

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Failed to generate token"))
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"user": fiber.Map{
				"uuid":     user.UUID,
				"id":       user.ID,
				"name":     user.Name,
				"nickname": user.Nickname,
				"picture":  user.Picture,
			},
			"token": token,
		},
	})
}
