package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID     uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey;not null;unique" json:"uuid"`
	ID       int64     `gorm:"column:id;not null;unique" json:"id"`
	Name     string    `gorm:"column:name;type:text;not null" json:"name"`
	Nickname string    `gorm:"column:nickname;type:text;not null" json:"nickname"`
	Picture  string    `gorm:"column:picture;type:text;not null" json:"picture"`
}
