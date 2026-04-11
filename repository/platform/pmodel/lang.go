package pmodel

import "gorm.io/gorm"

type Lang struct {
	gorm.Model

	Code    int
	Lang    string
	Content string
}
