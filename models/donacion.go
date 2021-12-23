package models

import (
	"gorm.io/gorm"
)

type Donacion struct {
	gorm.Model
	Nombre              string             `json:"nombre"`
	ParroquiaID         uint               `json:"id_parroquia"`
	Imagen              string             `json:"imagen_perfil"`
	ImagenReserva       string             `json:"imagen_banner"`
	Meta                float64            `json:"meta"`
	Descripcion         string             `json:"descripcion"`
	CategoriaDonacionID uint               `json:"id_categoria_donacion"`
	CategoriaDonacion   *CategoriaDonacion `json:"categoria_donacion"`
	Parroquia           *Parroquia         `json:"parroquia"`
	Aportaciones        []*Aportacion      `json:"aportaciones"`
}

func (Donacion) TableName() string {
	return "donacion"
}
