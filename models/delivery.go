package models

import "gorm.io/gorm"

type Delivery struct {
	gorm.Model
	Nombre string `json:"nombre"`
}
