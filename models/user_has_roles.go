package models

import "gorm.io/gorm"

type UserHasRoles struct {
	gorm.Model
	ID     uint `json:"id" gorm:"primaryKey"`
	Role   uint `json:"role"`
	UserID uint `json:"user"`
}
