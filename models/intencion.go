package models

import (
	"gorm.io/gorm"
)

type Intencion struct {
	gorm.Model
	FielID       uint         `json:"id_fiel"`
	ParroquiaID  uint         `json:"id_parroquia"`
	MisaID       uint         `json:"id_misa"`
	Tipo         string       `json:"tipo"`
	Nombre       string       `json:"nombre"`
	Monto        float64      `json:"monto"`
	TokenTarjeta string       `json:"token_tarjeta" gorm:"-"`
	Transaccion  *Transaccion `json:"transaccion"  gorm:"polymorphic:TipoPago"`
	Parroquia    *Parroquia   `json:"parroquia"`
	Misa         *Misa        `json:"misa"`
	Fiel         *Fiel        `json:"fiel"`
}

func (Intencion) TableName() string {
	return "intencion"
}
