package models

import "gorm.io/gorm"

type EmprendimientoImagen struct {
	gorm.Model
	Imagen           string          `json:"imagen"`
	EmprendimientoID uint            `json:"id_emprendimiento"`
	Emprendimiento   *Emprendimiento `json:"emprendimiento"`
}
