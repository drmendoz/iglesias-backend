package models

import "gorm.io/gorm"

type MatrimonioImagen struct {
	gorm.Model
	Imagen       string      `json:"imagen"`
	MatrimonioID uint        `json:"id_publicacion"`
	Matrimonio   *Matrimonio `json:"publicacion,omitempty"`
}
