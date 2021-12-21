package models

import "gorm.io/gorm"

type Publicidad struct {
	gorm.Model
	Imagen      string     `json:"imagen"`
	Empresa     string     `json:"empresa"`
	Documento   string     `json:"documento"`
	Telefono    string     `json:"telefono"`
	Prioridad   int        `json:"prioridad" gorm:"default:1"`
	ParroquiaID uint       `json:"etapa_id"`
	Parroquia   *Parroquia `json:"parroquia"`
}
