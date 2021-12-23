package models

import "gorm.io/gorm"

type CategoriaDonacion struct {
	gorm.Model
	Nombre     string      `json:"nombre"`
	Imagen     string      `json:"imagen"`
	Donaciones []*Donacion `json:"donaciones,omitempty"`
}

func (CategoriaDonacion) TableName() string {
	return "categoria_donacion"
}
