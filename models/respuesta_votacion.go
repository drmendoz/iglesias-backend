package models

import "gorm.io/gorm"

type RespuestaVotacion struct {
	gorm.Model
	OpcionVotacionID uint           `json:"id_opcion_votacion"`
	OpcionVotacion   OpcionVotacion `json:"opcion_votacion"`
	FielID           uint           `json:"id_residente"`
	Fiel             Fiel           `json:"residente"`
}

func (RespuestaVotacion) TableName() string {
	return "respuesta_votacion"
}
