package models

import (
	"gorm.io/gorm"
)

type ModuloEtapa struct {
	gorm.Model
	Modulo  string `json:"modulo"`
	Imagen  string `json:"imagen"`
	Estado  bool   `json:"estado"`
	EtapaID uint   `json:"id_etapa"`
	Etapa   *Etapa `json:"etapa"`
}
