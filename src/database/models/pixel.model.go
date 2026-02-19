package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pixel struct {
	gorm.Model
	X             uint16    `gorm:"column:x;not null;primaryKey" json:"x"`
	Y             uint16    `gorm:"column:y;not null;primaryKey" json:"y"`
	Color         uint8     `gorm:"column:color;not null" json:"c"`
	User          uuid.UUID `gorm:"column:user;type:uuid" json:"u"`
	UpdatedByUser *User     `gorm:"foreignKey:User;references:UUID"`
}
