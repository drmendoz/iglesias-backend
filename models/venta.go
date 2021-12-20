package models

import "gorm.io/gorm"

type Venta struct {
	gorm.Model
	Imagen  string `json:"imagen"`
	Nombre  string `json:"nombre"`
	Empresa string `json:"empresa"`
}

func (Venta) TableName() string {
	return "venta"
}
