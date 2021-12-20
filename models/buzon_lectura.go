package models

import "gorm.io/gorm"

type BuzonLectura struct {
	gorm.Model
	BuzonID   uint    `json:"id_buzon"`
	Buzon     Buzon   `json:"buzon"`
	UsuarioID uint    `json:"id_usuario"`
	Usuario   Usuario `json:"usuario"`
}

func (BuzonLectura) TableName() string {
	return "buzon_lectura"
}
