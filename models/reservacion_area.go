package models

import (
	"time"

	"gorm.io/gorm"
)

type ReservacionAreaSocial struct {
	gorm.Model
	HoraInicio     time.Time    `json:"hora_inicio"`
	HoraFin        time.Time    `json:"hora_fin"`
	ValorCancelado float64      `json:"valor_cancelado"`
	FielID         uint         `json:"id_residente"`
	AreaSocialID   uint         `json:"id_area_social"`
	TransaccionID  *uint        `json:"id_transaccion,omitempty"`
	Transaccion    *Transaccion `json:"transaccion,omitempty"`
	AreaSocial     *AreaSocial  `json:"area_social"`
	Fiel           *Fiel        `json:"residente,omitempty"`
}
