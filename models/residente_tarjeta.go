package models

import "gorm.io/gorm"

type FielTarjeta struct {
	gorm.Model
	TokenTarjeta  string         `json:"token"`
	FielID        uint           `json:"id_fiel"`
	Fiel          *Fiel          `json:"fiel"`
	Transacciones []*Transaccion `json:"transacciones"`
}
