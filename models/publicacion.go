package models

import "gorm.io/gorm"

type Publicacion struct {
	gorm.Model
	Titulo              string               `json:"titulo"`
	Descripcion         string               `json:"descripcion"`
	FielID              uint                 `json:"id_fiel"`
	ParroquiaID         uint                 `json:"id_parroquia"`
	Parroquia           *Parroquia           `json:"parroquia"`
	Fiel                *Fiel                `json:"fiel"`
	IsVertical          bool                 `json:"is_vertical" gorm:"default:false"`
	ImagenesPublicacion []*PublicacionImagen `json:"imagenes"`
}

func (Publicacion) TableName() string {
	return "publicacion"
}
