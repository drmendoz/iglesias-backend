package models

import "gorm.io/gorm"

type Autorizado struct {
	gorm.Model
	Conductor      string          `json:"conductor"`
	RazonSocial    string          `json:"razon_social"`
	Cedula         string          `json:"cedula"`
	Imagen         string          `json:"imagen"`
	Correo         string          `json:"correo"`
	Telefono       string          `json:"telefono"`
	Ano            string          `json:"ano"`
	Placa          string          `json:"placa"`
	Ruc            string          `json:"ruc"`
	Marca          string          `json:"marca"`
	Modelo         string          `json:"modelo"`
	TipoUsuario    string          `json:"tipo_usuario" gorm:"default:'VISITA'; gorm:type:enum('RESIDENTE','EMPLEADO','EXPRESO', 'FAMILIAR')"`
	Pdf            string          `json:"pdf"`
	PublicadorID   uint            `json:"id_publicador"`
	EtapaID        uint            `json:"id_etapa"`
	Autorizaciones []*Autorizacion `json:"autorizaciones"`
	Etapa          *Etapa          `json:"etapa,omitempty"`
	Publicador     *Usuario        `json:"publicador,omitempty"`
}

func (Autorizado) TableName() string {
	return "autorizado"
}
