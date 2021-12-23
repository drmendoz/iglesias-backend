package models

import "gorm.io/gorm"

type HorarioEntrada struct {
	gorm.Model
	Descripcion string
	HorarioID   uint
	Horario     Horario
}

func (HorarioEntrada) TableName() string {
	return "horario_entrada"
}
