package models

import "gorm.io/gorm"

type ModulosParroquia struct {
	gorm.Model
	ParroquiaID    uint       `json:"id_parroquia"`
	Parroquia      *Parroquia `json:"parroquia"`
	Horario        bool       `json:"horario" gorm:"default:false"`
	Actividad      bool       `json:"actividad" gorm:"default:false"`
	Emprendimiento bool       `json:"emprendimiento" gorm:"default:false"`
	Intencion      bool       `json:"intencion" gorm:"default:false"`
	Musica         bool       `json:"musica" gorm:"default:false"`
	Ayudemos       bool       `json:"ayudemos" gorm:"default:false"`
	Misa           bool       `json:"misa" gorm:"default:false"`
	Curso          bool       `json:"curso" gorm:"default:false"`
} // Agregar matrimonio, publicacion y galería

func (ModulosParroquia) TableName() string {
	return "modulos_parroquia"
}
