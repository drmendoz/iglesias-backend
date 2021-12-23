package models

import (
	"time"

	"gorm.io/gorm"
)

type Actividad struct {
	gorm.Model
	Titulo      string     `json:"titulo"`
	Descripcion string     `json:"descripcion"`
	Imagen      string     `json:"imagen"`
	Video       string     `json:"video"`
	FechaInicio *time.Time `json:"fecha_inicio"`
	TieneLimite bool       `json:"tiene_limite"`
	Cupo        int        `json:"cupo"`
	ParroquiaID uint       `json:"id_parroquia"`
	Parroquia   *Parroquia `json:"parroquia"`
}
