package models

import "gorm.io/gorm"

type Administrativo struct {
	gorm.Model
	Nombre      string `json:"nombre"`
	Imagen      string `json:"imagen"`
	Cedula      string `json:"cedula"`
	Telefono    string `json:"telefono"`
	Celular     string `json:"celular"`
	Correo      string `json:"correo"`
	Cargo       string `json:"cargo"`
	ParroquiaID uint   `json:"id_etapa"`
	Etapa       *Etapa `json:"etapa,omitempty"`
}
