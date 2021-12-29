package models

import (
	"time"

	"gorm.io/gorm"
)

type Matrimonio struct {
	gorm.Model
	Familia string `json:"familia"`

	FechaMatrimonio *time.Time    `json:"fecha_matrimonio"`
	Estado          string        `json:"estado" `
	DiferenciaFecha time.Duration `json:"diferencia_fecha" gorm:"-"`
	FielID          uint          `json:"id_fiel"`
	Fiel            *Fiel         `json:"fiel"`
	ParroquiaID     uint          `json:"id_parroquia"`
	Parroquia       *Parroquia    `json:"parroquia"`
	Monto           float64       `json:"monto"`
	TokenTarjeta    string        `json:"token_tarjeta" gorm:"-"`
	Transaccion     *Transaccion  `json:"transaccion"  gorm:"polymorphic:TipoPago"`
	Imagen          string        `json:"imagen"`
}
