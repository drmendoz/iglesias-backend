package models

import "gorm.io/gorm"

type AdminMaster struct {
	gorm.Model
	UsuarioID  uint               `json:"id_usuario"`
	Token      string             `json:"token,omitempty" gorm:"-" `
	ContraHash string             `json:"-" gorm:"-"`
	Usuario    *Usuario           `json:"usuario"`
	EsMaster   bool               `json:"is_master" gorm:"default:false"`
	Permisos   AdminMasterPermiso `json:"permisos"`
}

func (AdminMaster) TableName() string {
	return "admin_master"
}
