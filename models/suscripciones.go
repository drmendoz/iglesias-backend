package models

import (
	"time"

	"gorm.io/gorm"
)

type Suscripcion struct {
	gorm.Model
	FechaInicio time.Time  `json:"fecha_inicio"`
	FechaFin    time.Time  `json:"fecha_fin"`
	Continua    bool       `json:"continua"`
	ResidenteID uint       `json:"id_residente"`
	Residente   *Residente `json:"residente,omitempty"`
}

func (Suscripcion) TableName() string {
	return "suscripcion"
}
