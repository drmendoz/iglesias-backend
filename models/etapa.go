package models

import (
	"gorm.io/gorm"
)

type Etapa struct {
	gorm.Model
	Nombre                     string        `json:"nombre"`
	Latitud                    float64       `json:"lat"`
	Longitud                   float64       `json:"lng"`
	Direccion                  string        `json:"direccion"`
	Correo                     string        `json:"correo"`
	Telefono                   string        `json:"telefono"`
	Imagen                     string        `json:"imagen"`
	NombreUrbanizacion         string        `json:"nombre_urbanizacion,omitempty" gorm:"-"`
	UrbanizacionID             uint          `json:"id_urbanizacion"`
	FechaAlicuota              int           `json:"fecha_alicuota" gorm:"default:1"`
	ValorAlicuota              float64       `json:"valor_alicuota" gorm:"default:0"`
	NombreBanco                string        `json:"nombre_banco"`
	TipoCuenta                 string        `json:"tipo_cuenta"`
	NumeroCuenta               string        `json:"numero_cuenta"`
	TipoDocumento              string        `json:"tipo_documento"`
	NumeroDocumento            string        `json:"numero_documento"`
	Casa                       []Casa        `json:"-"`
	PagosTarjeta               bool          `json:"pagos_tarjeta" gorm:"default:false; gorm:type:boolean; column:pagos_tarjeta" mapstructure:"pagos_tarjeta"`
	ModuloMarket               bool          `json:"modulo_market" gorm:"default:true; gorm:type:boolean; column:modulo_market" mapstructure:"modulo_market"`
	ModuloPublicacion          bool          `json:"modulo_publicacion" gorm:"default:true; gorm:type:boolean; column:modulo_publicacion" mapstructure:"modulo_publicacion"`
	ModuloVotacion             bool          `json:"modulo_votacion" gorm:"default:true; gorm:type:boolean; column:modulo_votacion" mapstructure:"modulo_votacion"`
	ModuloAreaSocial           bool          `json:"modulo_area_social" gorm:"default:true; gorm:type:boolean; column:modulo_area_social" mapstructure:"modulo_area_social"`
	ModuloEquipoAdministrativo bool          `json:"modulo_equipo" gorm:"default:true; gorm:type:boolean; column:modulo_equipo" mapstructure:"modulo_equipo"`
	ModuloHistoria             bool          `json:"modulo_historia" gorm:"default:true; gorm:type:boolean; column:modulo_historia" mapstructure:"modulo_historia"`
	ModuloBitacora             bool          `json:"modulo_bitacora" gorm:"default:true; gorm:type:boolean; column:modulo_bitacora" mapstructure:"modulo_bitacora"`
	Urbanizacion               *Urbanizacion `json:"urbanizacion,omitempty"`
	FormularioEntrada          bool          `json:"formulario_entrada" gorm:"default:false; gorm:type:boolean; column:formulario_entrada" mapstructure:"formulario_entrada"`
	FormularioSalida           bool          `json:"formulario_salida" gorm:"default:false; gorm:type:boolean; column:formulario_salida" mapstructure:"formulario_salida"`
	ModuloAlicuota             bool          `json:"modulo_alicuota" gorm:"default:true; gorm:type:boolean; column:modulo_alicuota" mapstructure:"modulo_alicuota"`
	ModuloEmprendimiento       bool          `json:"modulo_emprendimiento" gorm:"default:true; gorm:type:boolean; column:modulo_emprendimiento" mapstructure:"modulo_emprendimiento"`
	ModuloCamaras              bool          `json:"modulo_camaras" gorm:"default:true; gorm:type:boolean; column:modulo_camaras" mapstructure:"modulo_camaras"`
	ModuloDirectiva            bool          `json:"modulo_directiva" gorm:"default:true; gorm:type:boolean; column:modulo_directiva" mapstructure:"modulo_directiva"`
	ModuloGaleria              bool          `json:"modulo_galeria" gorm:"default:true; gorm:type:boolean; column:modulo_galeria" mapstructure:"modulo_galeria"`
	ModuloHorarios             bool          `json:"modulo_horarios" gorm:"default:true; gorm:type:boolean; column:modulo_horarios" mapstructure:"modulo_horarios"`
	ModuloMiRegistro           bool          `json:"modulo_mi_registro" gorm:"default:true; gorm:type:boolean; column:modulo_mi_registro" mapstructure:"modulo_mi_registro"`
	ModuloAutorizacion         bool          `json:"modulo_autorizacion" gorm:"default:true; gorm:type:boolean; column:modulo_autorizacion" mapstructure:"modulo_autorizacion"`
}

func (Etapa) TableName() string {
	return "etapa"
}
