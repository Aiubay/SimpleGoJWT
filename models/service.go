package models

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	ID           uint   `json:"id"`
	Sr_number    string `json:"name"`
	Sr_status    string `json:"email"`
	On_behalf_of string `json:"on_behalf_of"`
	CreatedBy    uint   `json:"created_by"`
	UpdatedBy    uint   `json:"updated_by"`
}
