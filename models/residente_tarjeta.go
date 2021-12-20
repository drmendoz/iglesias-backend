package models

import "gorm.io/gorm"

type ResidenteTarjeta struct {
	gorm.Model
	TokenTarjeta  string         `json:"token"`
	ResidenteID   uint           `json:"id_residente"`
	Residente     *Residente     `json:"residente"`
	Transacciones []*Transaccion `json:"transacciones"`
}
