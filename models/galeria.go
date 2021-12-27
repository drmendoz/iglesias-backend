package models

import "gorm.io/gorm"

type Galeria struct {
	gorm.Model
	Nombre      string     `json:"nombre"`
	Imagen      string     `json:"imagen"`
	ParroquiaID uint       `json:"id_parroquia"`
	Parroquia   *Parroquia `json:"parroquia"`
}
