package models

import "gorm.io/gorm"

type PermisoUsuario struct {
	gorm.Model
	PermisoID uint    `json:"id_permiso"`
	Permiso   Permiso `json:"permiso"`
	UsuarioID uint    `json:"id_usuario"`
	Usuario   Usuario `json:"usuario"`
}

func (PermisoUsuario) TableName() string {
	return "permiso_usuario"
}
