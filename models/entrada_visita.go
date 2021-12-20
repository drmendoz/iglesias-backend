package models

import (
	"time"

	"gorm.io/gorm"
)

type EntradaVisita struct {
	gorm.Model
	HoraEntrada time.Time `json:"hora_entrada"`
	Placa       string    `json:"placa"`
	TipoEntrada string    `json:"tipo_entrada" gorm:"default:'VEHICULO'; gorm:type:enum('VEHICULO','MOTO','CAMINANDO', 'BICICLETA')"`
}

func (EntradaVisita) TableName() string {
	return "entrada"
}
