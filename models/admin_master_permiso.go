package models

import "gorm.io/gorm"

type AdminMasterPermiso struct {
	gorm.Model
	AdminMasterID uint         `json:"id_admin_master"`
	AdminMaster   *AdminMaster `json:"admin_master"`
	Iglesia       bool         `json:"iglesia" gorm:"default:false"`
	Parroquia     bool         `json:"parroquia" gorm:"default:false"`
	Administrador bool         `json:"administrador" gorm:"default:false"`
	Modulo        bool         `json:"modulo" gorm:"default:false"`
	Recuadacion   bool         `json:"recaudacion" gorm:"default:false"`
	Usuario       bool         `json:"usuario" gorm:"default:false"`
	Categoria     bool         `json:"categoria"`
}

func (AdminMasterPermiso) TableName() string {
	return "admin_master_permiso"
}
