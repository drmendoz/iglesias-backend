package models

import "gorm.io/gorm"

type AdminEtapa struct {
	gorm.Model
	UsuarioID          uint              `json:"id_usuario"`
	Cedula             string            `json:"cedula"`
	NombreEtapa        *string           `json:"nombre_etapa,omitempty" gorm:"-"`
	NombreUrbanizacion *string           `json:"nombre_urbanizacion,omitempty" gorm:"-"`
	Token              string            `json:"token,omitempty" gorm:"-"`
	ContraHash         string            `json:"-" gorm:"-"`
	EsMaster           bool              `json:"is_master" gorm:"default:false"`
	Modulos            *Modulos          `json:"modulos,omitempty" gorm:"-"`
	EtapaID            uint              `json:"id_etapa"`
	Usuario            *Usuario          `json:"usuario"`
	Etapa              *Etapa            `json:"etapa"`
	Pdf                string            `json:"documento"`
	Permisos           AdminEtapaPermiso `json:"permisos"`
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

func (AdminEtapa) TableName() string {
	return "admin_etapa"
}
