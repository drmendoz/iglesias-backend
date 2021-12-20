package models

import "gorm.io/gorm"

type LecturaHistoria struct {
	gorm.Model
	FielID         uint `json:"id_residente"`
	HistoriaFielID uint `json:"id_historia"`
	Fiel           *Fiel
	HistoriaFiel   *HistoriaFiel
}

func (LecturaHistoria) TableName() string {
	return "lectura_historia"
}
