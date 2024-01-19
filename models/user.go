package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint         `json:"id" gorm:"primaryKey"`
	Name     string       `json:"name"`
	Email    string       `json:"email" gorm:"unique"`
	Password []byte       `json:"password"`
	Role     UserHasRoles `json:"role" gorm:"foreignKey:UserID;association_foreignkey:ID"`
}

type PayloadLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PayloadRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
