package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pixel struct {
	X             int        `gorm:"column:x;not null;primaryKey" json:"x"`
	Y             int        `gorm:"column:y;not null;primaryKey" json:"y"`
	Color         int        `gorm:"column:color;not null" json:"c"`
	UpdatedBy     *uuid.UUID `gorm:"column:updated_by;type:uuid" json:"u"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBy;references:UUID"`
}

type PixelModel struct {
	gorm.Model
	Pixel
}
