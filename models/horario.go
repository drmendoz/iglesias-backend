package models

import "gorm.io/gorm"

type Horario struct {
	gorm.Model
	Nombre           string            `json:"nombre"`
	Imagen           string            `json:"imagen"`
	Telefono         string            `json:"telefono"`
	ParroquiaID      uint              `json:"id_parroquia"`
	Parroquia        *Parroquia        `json:"parroquia"`
	Horarios         []string          `json:"horarios" gorm:"-"`
	HorariosEntradas []*HorarioEntrada `json:"-" `
}

func (Horario) TableName() string {
	return "horario"
}
