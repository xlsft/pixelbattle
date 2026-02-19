package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUID uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey;not null;unique" json:"uuid"`
	TGID int64     `gorm:"column:id;not null;unique" json:"id"`
	Name string    `gorm:"column:name;type:text;not null" json:"name"`
}

type UserModel struct {
	gorm.Model
	UUID uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey;not null;unique" json:"uuid"`
	ID   int64     `gorm:"column:id;not null;unique" json:"id"`
	Name string    `gorm:"column:name;type:text;not null" json:"name"`
}
