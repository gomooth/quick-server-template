package pmodel

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model

	Genre  uint8
	UserID uint
}
