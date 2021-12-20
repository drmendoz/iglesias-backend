package models

import "gorm.io/gorm"

type ImagenGaleria struct {
	gorm.Model
	Titulo  string `json:"titulo"`
	Imagen  string `json:"imagen"`
	EtapaID uint   `json:"id_etapa"`
	Etapa   *Etapa `json:"etapa"`
}

func (ImagenGaleria) TableName() string {
	return "galeria"
}
