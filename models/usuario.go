package models

import "gorm.io/gorm"

type Usuario struct {
	gorm.Model
	Nombre          string  `json:"nombres,omitempty"`
	Apellido        string  `json:"apellido"`
	Telefono        string  `json:"telefono,omitempty"`
	Correo          string  `json:"correo" `
	Celular         string  `json:"celular"`
	Usuario         string  `json:"usuario,omitempty"`
	Contrasena      string  `json:"contrasena,omitempty" `
	RandomNumToken  string  `json:"-"`
	Imagen          string  `json:"imagen"`
	CodigoTemporal  string  `json:"codigo_temporal,omitempty"`
	Cedula          *string `json:"cedula"`
	Fiel            *Fiel   `json:"residente,omitempty"`
	ViejaContrasena string  `json:"vieja_contrasena,omitempty" gorm:"-"`
}

func (Usuario) TableName() string {
	return "usuario"
}
