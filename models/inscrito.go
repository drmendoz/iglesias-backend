package models

import "gorm.io/gorm"

type Inscrito struct {
	gorm.Model
	Nombre       string       `json:"nombre"`
	Apellido     string       `json:"apellido"`
	TokenTarjeta string       `json:"token_tarjeta" gorm:"-"`
	Monto        float64      `json:"monto" gorm:"-"`
	CursoID      uint         `json:"id_curso"`
	Curso        *Curso       `json:"curso"`
	FielID       uint         `json:"id_fiel"`
	Fiel         *Fiel        `json:"fiel"`
	Transaccion  *Transaccion `json:"transaccion" gorm:"polymorphic:TipoPago"`
}
