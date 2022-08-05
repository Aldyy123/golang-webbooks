package Model

import (
	"gorm.io/gorm"
)

type Books struct {
	gorm.Model
	ID          int `gorm:"primaryKey;autoIncrement:true"`
	Name        string
	Author      string
	Description string
	ImageCover  string
}
