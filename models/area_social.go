package models

import (
	"gorm.io/gorm"
)

type AreaSocial struct {
	gorm.Model
	Nombre                   string                   `json:"nombre"`
	ParroquiaID              uint                     `json:"id_etapa"`
	Imagen                   string                   `json:"imagen"`
	ImagenReserva            string                   `json:"imagen_reserva"`
	IsPublica                bool                     `json:"is_publica" gorm:"default:false"`
	Precio                   float64                  `json:"precio"`
	SeleccionCosto           string                   `json:"seleccionCosto"`
	Normas                   string                   `json:"normas"`
	TipoAforo                string                   `json:"tipoAforo"`
	TipoArea                 string                   `json:"tipoArea"`
	Aforo                    int                      `json:"aforo"`
	Estado                   string                   `json:"estado" gorm:"-"`
	TiempoReservacionMinutos int                      `json:"tiempo_reservacion_minutos"`
	ReservasFielMes          int                      `json:"reservas_mes_residente"`
	Horario                  []*AreaHorario           `json:"horario,omitempty" gorm:"-"`
	Horarios                 []AreaHorario            `json:"horarios,omitempty"`
	Reservaciones            *[]ReservacionAreaSocial `json:"reservaciones,omitempty"`
	Etapa                    *Etapa                   `json:"etapa,omitempty"`
}

func (AreaSocial) TableName() string {
	return "area_social"
}
