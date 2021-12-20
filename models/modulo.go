package models

import "gorm.io/gorm"

type Modulo struct {
	gorm.Model
	Nombre string `json:"nombre"`
}

func (Modulo) TableName() string {
	return "modulo"
}
