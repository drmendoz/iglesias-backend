package models

import "gorm.io/gorm"

type Musica struct {
	gorm.Model
	Titulo      string     `json:"titulo"`
	Compositor  string     `json:"compositor"`
	Telefono    string     `json:"telefono"`
	Media       string     `json:"media"`
	FielID      uint       `json:"id_fiel"`
	Fiel        *Fiel      `json:"fiel"`
	ParroquiaID uint       `json:"id_parroquia"`
	Parroquia   *Parroquia `json:"parroquia"`
	Estado      string     `json:"estado" gorm:"type:enum('PEN','PUB')"`
}
