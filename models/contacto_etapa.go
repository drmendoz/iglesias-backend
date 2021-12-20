package models

import "gorm.io/gorm"

type ContactoEtapa struct {
	gorm.Model
	Nombre   string `json:"contacto"`
	Imagen   string `json:"imagen"`
	Telefono string `json:"telefono"`
	Horario  string `json:"horario"`
	EtapaID  uint   `json:"id_etapa"`
	Etapa    *Etapa `json:"etapa"`
}
