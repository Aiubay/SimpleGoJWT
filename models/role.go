package models

type Role struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Role        string `json:"role"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}
