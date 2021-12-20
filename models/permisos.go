package models

import "gorm.io/gorm"

type Permiso struct {
	gorm.Model
	Nombre string `json:"nombre"`
}

func (Permiso) TableName() string {
	return "permiso"
}
