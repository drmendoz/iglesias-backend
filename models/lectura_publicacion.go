package models

import "gorm.io/gorm"

type LecturaPublicacion struct {
	gorm.Model
	PublicacionID uint
	UsuarioID     uint
	Usuario       *Usuario
	Publicacion   *Publicacion
}

func (LecturaPublicacion) TableName() string {
	return "lectura_publicacion"
}
