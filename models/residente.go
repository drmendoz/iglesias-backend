package models

import (
	"time"

	"gorm.io/gorm"
)

type Residente struct {
	gorm.Model
	UsuarioID                    uint                `json:"id_usuario"`
	CasaID                       uint                `json:"id_casa"`
	TokenNotificacion            string              `json:"token_notificacion,omitempty"`
	Cedula                       string              `json:"cedula"`
	Token                        string              `json:"token,omitempty" gorm:"-"`
	FechaNacimiento              string              `json:"fecha_nacimiento,omitempty"`
	Confirmacion                 bool                `json:"confirmacion" gorm:"default:true"`
	IsPrincipal                  bool                `json:"is_principal" gorm:"default:false"`
	Autorizacion                 bool                `json:"autorizacion" gorm:"default:false"`
	Usuario                      *Usuario            `json:"usuario,omitempty"`
	ContraHash                   string              `json:"-" gorm:"-"`
	Casa                         *Casa               `json:"casa,omitempty"`
	Pdf                          string              `json:"documento"`
	Mensaje                      string              `json:"mensaje" gorm:"-"`
	Emprendimientos              []*Emprendimiento   `json:"emprendimientos,omitempty"`
	Tarjetas                     []*ResidenteTarjeta `json:"tarjetas,omitempty"`
	SesionIniciada               bool                `json:"sesion_iniciada"`
	VisualizacionEmprendimiento  *time.Time          `json:"-" `
	VisualizacionBitacora        *time.Time          `json:"-" `
	VisualizacionBuzon           *time.Time          `json:"-"`
	VisualizacionGaleria         *time.Time          `json:"-"`
	VisualizacionVotacion        *time.Time          `json:"-"`
	VisualizacionCamara          *time.Time          `json:"-" `
	VisualizacionAlicuota        *time.Time          `json:"-" `
	VisualizacionAreaSocial      *time.Time          `json:"-" `
	VisualizacionAdministradores *time.Time          `json:"-" `
	VisualizacionReservas        *time.Time          `json:"-" `
}

func (Residente) TableName() string {
	return "residente"
}
