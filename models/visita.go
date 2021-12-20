package models

import (
	"time"

	"gorm.io/gorm"
)

type Visita struct {
	gorm.Model
	Nombre              string         `json:"nombres"`
	Apellidos           string         `json:"apellidos"`
	Cedula              string         `json:"cedula"`
	Placa               string         `json:"placa"`
	Imagen              string         `json:"imagen"`
	Imagen2             string         `json:"imagen2"`
	Motivo              string         `json:"motivo"`
	CasaID              uint           `json:"id_casa"`
	Casa                *Casa          `json:"casa,omitempty"`
	Estado              string         `gorm:"default:'PENDIENTE'; gorm:type:enum('PENDIENTE','ACEPTADA','RECHAZADA', 'ESPERANDO', 'SIN_RESPUESTA')"`
	RespuestaPorLlamada bool           `json:"respuesta_por_llamada" gorm:"default:false"`
	TipoEntrada         string         `json:"tipo_entrada" gorm:"default:'VEHICULO'; gorm:type:enum('VEHICULO','MOTO','CAMINANDO', 'BICICLETA')"`
	Nuevo               bool           `json:"nuevo" gorm:"default:true"`
	PublicadorID        uint           `json:"id_guardia"`
	UsuarioID           uint           `json:"id_usuario"`
	ParroquiaID         uint           `json:"id_etapa"`
	DiaCreacion         time.Time      `json:"dia_creacion"`
	EntradaID           *uint          `json:"id_entrada"`
	SalidaID            *uint          `json:"id_salida"`
	AutorizacionID      *uint          `json:"id_autorizacion"`
	Entrada             *EntradaVisita `json:"entrada,omitempty"`
	Salida              *SalidaVisita  `json:"salida,omitempty"`
	TipoUsuario         string         `json:"tipo_usuario" gorm:"default:'VISITA'; gorm:type:enum('RESIDENTE','EMPLEADO','VISITA','EXPRESO','CONDUCTOR','EXPRESO','DELIVERY','FAMILIAR','MUDANZA')"`
	Vista               bool           `json:"vista" gorm:"default:false"`
	Etapa               *Etapa         `json:"etapa,omitempty"`
	Publicador          *Usuario       `json:"guardia,omitempty" gorm:"foreignKey:PublicadorID"`
	Usuario             *Usuario       `json:"contestador"`
	Autorizacion        *Autorizacion  `json:"autorizacion"`
	SalidaFirst         bool           `json:"is_salida_first"`
	SegundosRestantes   time.Duration  `json:"seconds_remaining" gorm:"-"`
	SegundosTotal       int            `json:"seconds_total" gorm:"-"`
}

func (Visita) TableName() string {
	return "visita"
}
