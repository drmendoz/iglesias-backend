package models

import "gorm.io/gorm"

type Permiso struct {
	gorm.Model
	Permiso   bool     `json:"permiso"`
	Modulo    string   `json:"modulo"`
	UsuarioID uint     `json:"id_usuario"`
	Usuario   *Usuario `json:"usuario"`
}

func (Permiso) TableName() string {
	return "permiso"
}
