package models

import "gorm.io/gorm"

type AdminGarita struct {
	gorm.Model
	UsuarioID         uint       `json:"id_usuario"`
	ContraHash        string     `json:"-" gorm:"-"`
	Token             string     `json:"token,omitempty" gorm:"-"`
	EtapaLabel        *EtapaInfo `json:"data_etapa,omitempty" gorm:"-"`
	UrbanizacionLabel *UrbInfo   `json:"data_urb,omitempty" gorm:"-"`
	Pdf               string     `json:"documento"`

	ParroquiaID uint     `json:"id_etapa"`
	Usuario     *Usuario `json:"usuario"`
}

type EtapaInfo struct {
	EtapaNombre string `json:"nombre"`
	Imagen      string `json:"imagen"`
}
type UrbInfo struct {
	Nombre string `json:"nombre"`
	Imagen string `json:"imagen"`
}

func (AdminGarita) TableName() string {
	return "admin_garita"
}
