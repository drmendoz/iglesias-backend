package models

import (
	"time"

	"gorm.io/gorm"
)

type SalidaVisita struct {
	gorm.Model
	HoraSalida time.Time `json:"hora_salida"`
	Placa      string    `json:"placa"`
	TipoSalida string    `json:"tipo_salida" gorm:"default:'VEHICULO'; gorm:type:enum('VEHICULO','MOTO','CAMINANDO', 'BICICLETA')"`
}

func (SalidaVisita) TableName() string {
	return "salida"
}
