package models

import (
	"time"

	"gorm.io/gorm"
)

type Alicuota struct {
	gorm.Model
	Valor         float64      `json:"valor"`
	CasaID        uint         `json:"id_casa"`
	Tipo          string       `json:"tipo" gorm:"default:'COMUN'; gorm:type:enum('COMUN','SALDO','EXTRAORDINARIA)"`
	Estado        string       `json:"estado" gorm:"default:'PENDIENTE';type:enum('PAGADO','PENDIENTE','VENCIDO)"`
	TransaccionID *uint        `json:"id_transaccion,omitempty"`
	Transaccion   *Transaccion `json:"transaccion,omitempty"`
	FechaPago     *time.Time   `json:"fecha_pago" `
	MesPago       time.Time    `json:"mes_pago" `
	Casa          *Casa        `json:"casa"`
}

func (Alicuota) TableName() string {
	return "alicuota"
}
