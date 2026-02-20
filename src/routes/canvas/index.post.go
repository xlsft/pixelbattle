package canvasRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
	"github.com/xlsft/pixelbattle/utils"
	"gorm.io/gorm"
)

type PixelRequest struct {
	X     uint16 `json:"x"`
	Y     uint16 `json:"y"`
	Color uint8  `json:"color"`
}

func HandlePost(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(models.User)
	db := database.UseDb()
	request := PixelRequest{}

	var pixel models.Pixel

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	if err := db.Model(&models.Pixel{}).Preload("UpdatedByUser").Where(&models.Pixel{X: request.X, Y: request.Y}).First(&pixel).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Database error"))
		}

		pixel = models.Pixel{
			User:  user.UUID,
			X:     request.X,
			Y:     request.Y,
			Color: request.Color,
		}

		if err := db.Model(&models.Pixel{}).Create(&pixel).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Failed to draw a pixel"))
		}
	} else {
		if tx := db.Model(&models.Pixel{}).Where(&pixel).Updates(models.Pixel{Color: request.Color, User: user.UUID}); tx.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(utils.DefineError("Failed to update a pixel"))
		}
	}

	PushEvents([]PixelRequest{{
		X:     pixel.X,
		Y:     pixel.Y,
		Color: pixel.Color,
	}})

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"user":    user.Name,
			"x":       pixel.X,
			"y":       pixel.Y,
			"color":   pixel.Color,
			"updated": pixel.UpdatedAt,
		},
	})
}
