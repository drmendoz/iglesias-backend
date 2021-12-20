package models

import (
	"time"

	"gorm.io/gorm"
)

type Votacion struct {
	gorm.Model
	Pregunta         string            `json:"pregunta"`
	EtapaID          uint              `json:"id_etapa"`
	UsuarioVotacion  bool              `json:"usuario_voto" gorm:"-"`
	FechaVencimiento time.Time         `json:"fecha_vencimiento"`
	Expiro           bool              `json:"expiro" gorm:"-"`
	TotalVotos       int               `json:"total_votos" gorm:"default:0"`
	Imagenes         string            `json:"imagenes_string"`
	ImagenesArray    []string          `json:"imagenes" gorm:"-"`
	Etapa            *Etapa            `json:"etapa"`
	Opciones         []*OpcionVotacion `json:"opciones"`
}

func (Votacion) TableName() string {
	return "votacion"
}
