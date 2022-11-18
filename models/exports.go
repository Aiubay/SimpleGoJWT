package models

import (
	"gorm.io/gorm"
)

type Excel struct {
	gorm.Model
	Id     int    `json:"id" gorm:"autoIncrement"`
	Random string `json:"random"`
	// Time   time.Time `json:"time"`
}
