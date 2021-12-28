package models

import "gorm.io/gorm"

type ModulosParroquia struct {
	gorm.Model
	ParroquiaID    uint       `json:"id_parroquia"`
	Parroquia      *Parroquia `json:"parroquia"`
	Horario        bool       `json:"horario" gorm:"default:true"`
	Actividad      bool       `json:"actividad" gorm:"default:true"`
	Emprendimiento bool       `json:"emprendimiento" gorm:"default:true"`
	Intencion      bool       `json:"intencion" gorm:"default:true"`
	Musica         bool       `json:"musica" gorm:"default:true"`
	Ayudemos       bool       `json:"ayudemos" gorm:"default:true"`
	Curso          bool       `json:"curso" gorm:"default:true"`
	Galeria        bool       `json:"galeria" gorm:"default:true"`
	Publicacion    bool       `json:"publicacion" gorm:"default:true"`
	Matrimonio     bool       `json:"matrimonio" gorm:"default:true"`
}

func (ModulosParroquia) TableName() string {
	return "modulos_parroquia"
}
