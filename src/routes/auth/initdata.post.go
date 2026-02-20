package authRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
	"github.com/xlsft/pixelbattle/utils"
)

func HandleInitDataPost(ctx *fiber.Ctx) error {
	request := utils.TelegramInitData{}
	db := database.UseDb()

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	if err := request.VerifyTelegramInitData(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	var user models.User
	data, err := request.ParseTelegramInitData()

	if err := db.Model(&models.User{}).Where("id = ?", data.User.ID).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Database error"))
		}

		user = models.User{
			UUID:     uuid.New(),
			ID:       data.User.ID,
			Name:     data.User.FirstName + " " + data.User.LastName,
			Nickname: data.User.Username,
			Picture:  data.User.PhotoURL,
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
