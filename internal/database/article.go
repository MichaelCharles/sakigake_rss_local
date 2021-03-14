package database

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model

	Title     string
	Author    string
	Content   string
	Published time.Time
}
