package models

import (
	"time"

	"gorm.io/gorm"
)

type AreaHorario struct {
	gorm.Model
	Dia                  string      `json:"dia"`
	HoraInicio           *time.Time  `json:"-"`
	HoraInicioSinFormato string      `json:"hora_inicio" gorm:"-"`
	HoraFin              *time.Time  `json:"-"`
	HoraFinSinFormato    string      `json:"hora_fin" gorm:"-"`
	FechaInicio          string      `json:"fecha_inicio"`
	FechaFin             string      `json:"fecha_fin"`
	AreaSocialID         uint        `json:"id_area"`
	AreaSocial           *AreaSocial `json:"area_social,omitempty"`
}
