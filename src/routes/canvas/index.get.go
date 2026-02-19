package canvasRoutes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/xlsft/pixelbattle/database"
	"github.com/xlsft/pixelbattle/database/models"
	"github.com/xlsft/pixelbattle/utils"
)

type PixelInfoQuery struct {
	X uint16 `query:"x"`
	Y uint16 `query:"y"`
}

func HandleGet(ctx *fiber.Ctx) error {
	db := database.UseDb()
	request := PixelRequest{}

	var pixel models.Pixel

	if err := ctx.QueryParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.DefineError(err.Error()))
	}

	if err := db.Model(&models.Pixel{}).Preload("UpdatedByUser").Where(&models.Pixel{X: request.X, Y: request.Y}).First(&pixel).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(utils.DefineError("Database error"))
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"user": fiber.Map{
				"name":    pixel.UpdatedByUser.Name,
				"updated": pixel.UpdatedByUser.UpdatedAt,
			},
			"x":     pixel.X,
			"y":     pixel.Y,
			"color": pixel.Color,
		},
	})
}
