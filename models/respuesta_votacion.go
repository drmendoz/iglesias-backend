package models

import "gorm.io/gorm"

type RespuestaVotacion struct {
	gorm.Model
	OpcionVotacionID uint           `json:"id_opcion_votacion"`
	OpcionVotacion   OpcionVotacion `json:"opcion_votacion"`
	ResidenteID      uint           `json:"id_residente"`
	Residente        Residente      `json:"residente"`
}

func (RespuestaVotacion) TableName() string {
	return "respuesta_votacion"
}
