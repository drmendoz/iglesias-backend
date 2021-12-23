package models

import "gorm.io/gorm"

type Transaccion struct {
	gorm.Model
	Estado             string       `json:"estado_pago"`
	DiaPago            string       `json:"dia_pago"`
	Monto              string       `json:"monto"`
	CodigoAutorizacion string       `json:"codigo_autorizacion"`
	Mensaje            string       `json:"mensaje"`
	Descripcion        string       `json:"descripcion"`
	TipoPagoID         uint         `json:"id_tipo_pago"`
	TipoPagoType       string       `json:"tipo_pago"`
	CategoriaID        uint         `json:"id_categoria"`
	EstadoDevolucion   string       `json:"estado_devolucion"`
	DetalleDevolucion  string       `json:"detalle_devolucion"`
	FielTarjetaID      uint         `json:"tarjeta_id"`
	FielTarjeta        *FielTarjeta `json:"tarjeta,omitempty"`
}
