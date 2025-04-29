package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	AccountID   string  `json:"account_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Budget      float64 `json:"budget"`
	Current     float64 `json:"current"`
}
