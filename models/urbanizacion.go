package models

import "gorm.io/gorm"

type Urbanizacion struct {
	gorm.Model
	Nombre    string   `json:"nombre"`
	Direccion string   `json:"direccion"`
	Telefono  string   `json:"telefono"`
	Correo    string   `json:"correo"`
	Ciudad    string   `json:"ciudad"`
	Latitud   float64  `json:"lat"`
	Longitud  float64  `json:"lng"`
	Imagen    string   `json:"imagen"`
	Etapas    []*Etapa `json:"etapas"`
}

func (Urbanizacion) TableName() string {
	return "urbanizacion"
}
