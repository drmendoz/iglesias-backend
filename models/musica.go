package models

import "gorm.io/gorm"

type Musica struct {
	gorm.Model
	Titulo      string     `json:"titulo"`
	Compositor  string     `json:"compositor"`
	Telefono    string     `json:"telefono"`
	Media       string     `json:"media"`
	ParroquiaID uint       `json:"id_parroquia"`
	Parroquia   *Parroquia `json:"parroquia"`
	Estado      string     `json:"estado" gorm:"type:enum('PEN','PUB')"`
}
