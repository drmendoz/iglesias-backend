package models

import (
	"gorm.io/gorm"
)

type Aportacion struct {
	gorm.Model
	FielID       uint         `json:"id_fiel"`
	DonacionID   uint         `json:"id_aportacion"`
	Fiel         *Fiel        `json:"fiel"`
	Donacion     *Donacion    `json:"donacion"`
	Mensaje      string       `json:"mensaje"`
	Monto        float64      `json:"monto"`
	TokenTarjeta string       `json:"token_tarjeta" gorm:"-"`
	Transaccion  *Transaccion `json:"transaccion"  gorm:"polymorphic:TipoPago"`
}

func (Aportacion) TableName() string {
	return "aportacion"
}
