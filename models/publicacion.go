package models

import "gorm.io/gorm"

type Publicacion struct {
	gorm.Model
	Titulo              string               `json:"titulo"`
	Cuerpo              string               `json:"cuerpo"`
	Leido               bool                 `json:"leido" gorm:"-"`
	Adjuntado           bool                 `json:"adjuntado" gorm:"default:false"`
	UsuarioID           uint                 `json:"id_usuario"`
	ParroquiaID         uint                 `json:"id_etapa"`
	Etapa               *Etapa               `json:"etapa"`
	Usuario             *Usuario             `json:"usuario"`
	IsVertical          bool                 `json:"is_vertical" gorm:"default:false"`
	ImagenesPublicacion []*PublicacionImagen `json:"imagenes"`
}

func (Publicacion) TableName() string {
	return "publicacion"
}
