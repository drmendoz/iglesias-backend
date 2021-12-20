package models

import (
	"time"

	"gorm.io/gorm"
)

type HistoriaResidente struct {
	gorm.Model
	Url         string     `json:"url"`
	IsVideo     bool       `json:"is_imagen"`
	FechaFin    time.Time  `json:"fecha_fin"`
	ResidenteID uint       `json:"id_autor"`
	Residente   *Residente `json:"autor,omitempty"`
}
