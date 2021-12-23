package models

import "gorm.io/gorm"

type Curso struct {
	gorm.Model
	ParroquiaID uint        `json:"id_parroquia"`
	Parroquia   *Parroquia  `json:"parroquia"`
	Titulo      string      `json:"titulo"`
	Descripcion string      `json:"descripcion"`
	Imagen      string      `json:"imagen"`
	Video       string      `json:"video"`
	Telefono    string      `json:"telefono"`
	FechaInicio string      `json:"fecha_inicio"`
	FechaFin    string      `json:"fecha_fin"`
	Precio      float64     `json:"precio"`
	TieneLimite bool        `json:"tiene_limite" `
	BotonPago   bool        `json:"boton_pago"`
	Cupo        int         `json:"cupo"`
	Inscritos   []*Inscrito `json:"inscritos"`
}
