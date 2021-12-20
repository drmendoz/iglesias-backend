package models

import "gorm.io/gorm"

type AdminParroquia struct {
	gorm.Model
	UsuarioID       uint                  `json:"id_usuario"`
	Cedula          string                `json:"cedula"`
	NombreParroquia *string               `json:"nombre_etapa,omitempty" gorm:"-"`
	Token           string                `json:"token,omitempty" gorm:"-"`
	ContraHash      string                `json:"-" gorm:"-"`
	EsMaster        bool                  `json:"is_master" gorm:"default:false"`
	Modulos         *Modulos              `json:"modulos,omitempty" gorm:"-"`
	ParroquiaID     uint                  `json:"id_parroquia"`
	Usuario         *Usuario              `json:"usuario"`
	Parroquia       *Parroquia            `json:"parroquia"`
	Permisos        AdminParroquiaPermiso `json:"permisos"`
}

type Modulos struct {
	ModuloMarket               bool `json:"market"`
	ModuloPublicacion          bool `json:"publicacion"`
	ModuloVotacion             bool `json:"votacion"`
	ModuloAreaSocial           bool `json:"area_social"`
	ModuloEquipoAdministrativo bool `json:"administradores"`
	ModuloHistoria             bool `json:"historias"`
	ModuloBitacora             bool `json:"bitacora"`
}

func (AdminParroquia) TableName() string {
	return "admin_parroquia"
}
