package models

import "gorm.io/gorm"

type FielTarjeta struct {
	gorm.Model
	TokenTarjeta  string         `json:"token"`
	FielID        uint           `json:"id_residente"`
	Fiel          *Fiel          `json:"residente"`
	Transacciones []*Transaccion `json:"transacciones"`
}
