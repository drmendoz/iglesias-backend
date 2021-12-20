package models

import (
	"gorm.io/gorm"
)

type Fiel struct {
	gorm.Model
	UsuarioID         uint           `json:"id_usuario"`
	TokenNotificacion string         `json:"token_notificacion,omitempty"`
	Cedula            string         `json:"cedula"`
	Token             string         `json:"token,omitempty" gorm:"-"`
	FechaNacimiento   string         `json:"fecha_nacimiento,omitempty"`
	Confirmacion      bool           `json:"confirmacion" gorm:"default:true"`
	ParroquiaID       *uint          `json:"id_parroquia"`
	Parroquia         *Parroquia     `json:"parroquia"`
	Usuario           *Usuario       `json:"usuario,omitempty"`
	ContraHash        string         `json:"-" gorm:"-"`
	Mensaje           string         `json:"mensaje" gorm:"-"`
	Tarjetas          []*FielTarjeta `json:"tarjetas,omitempty"`
	SesionIniciada    bool           `json:"sesion_iniciada"`
}

func (Fiel) TableName() string {
	return "fiel"
}
