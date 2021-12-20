package models

import (
	"time"

	"gorm.io/gorm"
)

type Suscripcion struct {
	gorm.Model
	FechaInicio time.Time `json:"fecha_inicio"`
	FechaFin    time.Time `json:"fecha_fin"`
	Continua    bool      `json:"continua"`
	FielID      uint      `json:"id_residente"`
	Fiel        *Fiel     `json:"residente,omitempty"`
}

func (Suscripcion) TableName() string {
	return "suscripcion"
}
