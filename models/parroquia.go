package models

import (
	"gorm.io/gorm"
)

type Parroquia struct {
	gorm.Model
	Nombre                  string           `json:"nombre"`
	Latitud                 float64          `json:"lat"`
	Longitud                float64          `json:"lng"`
	Direccion               string           `json:"direccion"`
	Correo                  string           `json:"correo"`
	Telefono                string           `json:"telefono"`
	Imagen                  string           `json:"imagen"`
	NombreBanco             string           `json:"nombre_banco"`
	TipoCuenta              string           `json:"tipo_cuenta"`
	NumeroCuenta            string           `json:"numero_cuenta"`
	TipoDocumento           string           `json:"tipo_documento"`
	NumeroDocumento         string           `json:"numero_documento"`
	BotonPagoIntencion      *bool            `json:"boton_pago_intencion"`
	BotonPagoCurso          *bool            `json:"boton_pago_curso"`
	BotonPagoEmprendimiento *bool            `json:"boton_pago_emprendimiento"`
	BotonPagoMatrimonio     *bool            `json:"boton_pago_matrimonio"`
	BotonPagoActividad      *bool            `json:"boton_pago_actividad"`
	BotonPagoMusica         *bool            `json:"boton_pago_musica"`
	BotonAyudemos           *bool            `json:"boton_pago_ayudemos" gorm:"default:false"`
	CostoEmprendimiento     float64          `json:"costo_emprendimiento"`
	IglesiaID               *uint            `json:"id_iglesia"`
	Iglesia                 *Iglesia         `json:"iglesia"`
	Modulos                 ModulosParroquia `json:"modulos"`
}

func (Parroquia) TableName() string {
	return "parroquia"
}
