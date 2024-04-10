package entities

// IMPORTANT!! Entities or better known as Domains

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ID   string `json:"id" binding:"required"`
	Date string `json:"date" binding:"required"`
}
