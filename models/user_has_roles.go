package models

import "gorm.io/gorm"

type UserHasRoles struct {
	gorm.Model
	ID     uint `json:"ID" gorm:"primaryKey"`
	Role   uint `json:"role"`
	UserID uint `json:"user"`
}
