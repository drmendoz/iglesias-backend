package models

import "gorm.io/gorm"

type OpcionVotacion struct {
	gorm.Model
	Descripcion string  `json:"opcion"`
	VotacionID  uint    `json:"id_votacion"`
	Conteo      int     `json:"total"`
	Color       string  `json:"color,omitempty" gorm:"-"`
	Porcentaje  float64 `json:"porcentaje" gorm:"-"`
}

func (OpcionVotacion) TableName() string {
	return "opcion_votacion"
}
