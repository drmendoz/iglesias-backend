package models

import "gorm.io/gorm"

type ParroquiaModulo struct {
	gorm.Model
	ParroquiaID uint       `json:"id_parroquia"`
	ModuloID    uint       `json:"id_modulo"`
	Parroquia   *Parroquia `json:"parroquia"`
	Modulo      *Modulo    `json:"modulo"`
}

func (ParroquiaModulo) TableName() string {
	return "parroquia_modulo"
}
