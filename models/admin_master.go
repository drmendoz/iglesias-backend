package models

import "gorm.io/gorm"

type AdminMaster struct {
	gorm.Model
	UsuarioID  uint               `json:"id_usuario"`
	Token      string             `json:"token,omitempty" gorm:"-" `
	ContraHash string             `json:"-" gorm:"-"`
	Usuario    *Usuario           `json:"usuario"`
	Permisos   AdminMasterPermiso `json:"permisos"`
	EsMaster   bool               `json:"is_master" gorm:"default:false"`
}

func (AdminMaster) TableName() string {
	return "admin_master"
}
