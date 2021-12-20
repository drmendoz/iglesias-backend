package models

import "gorm.io/gorm"

type LecturaHistoria struct {
	gorm.Model
	ResidenteID         uint `json:"id_residente"`
	HistoriaResidenteID uint `json:"id_historia"`
	Residente           *Residente
	HistoriaResidente   *HistoriaResidente
}

func (LecturaHistoria) TableName() string {
	return "lectura_historia"
}
