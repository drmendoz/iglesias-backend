package models

import "gorm.io/gorm"

type CategoriaMarket struct {
	gorm.Model
	Nombre          string            `json:"nombre"`
	Imagen          string            `json:"imagen"`
	Emprendimientos []*Emprendimiento `json:"emprendimientos,omitempty"`
}

func (CategoriaMarket) TableName() string {
	return "categoria_market"
}
