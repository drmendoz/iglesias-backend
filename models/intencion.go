package models

import (
	"gorm.io/gorm"
)

type Intencion struct {
	gorm.Model
	FielID        uint         `json:"id_fiel"`
	ParroquiaID   uint         `json:"id_parroquia"`
	MisaID        uint         `json:"id_misa"`
	TransaccionID uint         `json:"id_transaccion"`
	Tipo          string       `json:"tipo"`
	Nombre        string       `json:"nombre"`
	Transaccion   *Transaccion `json:"transaccion"`
	Parroquia     *Parroquia   `json:"parroquia"`
	Misa          *Misa        `json:"misa"`
	Fiel          *Fiel        `json:"fiel"`
}

func (Intencion) TableName() string {
	return "intencion"
}
