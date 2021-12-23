package models

import (
	"time"

	"gorm.io/gorm"
)

type Matrimonio struct {
	gorm.Model
	NombrePersona1     string              `json:"nombre_persona_1"`
	ApellidoPersona1   string              `json:"apellido_persona_1"`
	NombrePersona2     string              `json:"nombre_persona_2"`
	ApellidoPersona2   string              `json:"apellido_persona_2"`
	FechaMatrimonio    *time.Time          `json:"fecha_matrimonio"`
	FielID             uint                `json:"id_fiel"`
	Fiel               *Fiel               `json:"fiel"`
	ParroquiaID        uint                `json:"id_parroquia"`
	Parroquia          *Parroquia          `json:"parroquia"`
	Monto              float64             `json:"monto"`
	TokenTarjeta       string              `json:"token_tarjeta" gorm:"-"`
	Transaccion        *Transaccion        `json:"transaccion"  gorm:"polymorphic:TipoPago"`
	MatrimonioImagenes []*MatrimonioImagen `json:"-" `
	Imagenes           []string            `json:"imagenes" gorm:"-"`
}
