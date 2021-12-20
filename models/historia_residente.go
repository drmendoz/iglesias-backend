package models

import (
	"time"

	"gorm.io/gorm"
)

type HistoriaFiel struct {
	gorm.Model
	Url      string    `json:"url"`
	IsVideo  bool      `json:"is_imagen"`
	FechaFin time.Time `json:"fecha_fin"`
	FielID   uint      `json:"id_autor"`
	Fiel     *Fiel     `json:"autor,omitempty"`
}
