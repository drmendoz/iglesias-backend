package models

import "gorm.io/gorm"

type EtapaCamara struct {
	gorm.Model
	Nombre      string `json:"nombre"`
	Url         string `json:"url"`
	Imagen      string `json:"imagen"`
	ParroquiaID uint   `json:"id_etapa"`
	Etapa       *Etapa `json:"etapa,omitempty"`
}

func (EtapaCamara) TableName() string {
	return "etapa_camara"
}
