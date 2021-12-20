package models

import "gorm.io/gorm"

type PublicacionImagen struct {
	gorm.Model
	Imagen        string       `json:"imagen"`
	IsVideo       bool         `json:"is_video" gorm:"default:false"`
	PublicacionID uint         `json:"id_publicacion"`
	Publicacion   *Publicacion `json:"publicacion,omitempty"`
}
