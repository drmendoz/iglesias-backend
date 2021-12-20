package models

import (
	"time"

	"gorm.io/gorm"
)

type Misa struct {
	gorm.Model
	PadreID       uint       `json:"id_padre"`
	ParroquiaID   uint       `json:"id_parroquia"`
	CupoIntencion int        `json:"cupo_intencion"`
	FechaInicio   time.Time  `json:"fecha_inicio"`
	FechaFin      time.Time  `json:"fecha_fin"`
	Parroquia     *Parroquia `json:"parroquia"`
	Padre         *Padre     `json:"padre"`
}

func (Misa) TableName() string {
	return "misa"
}
