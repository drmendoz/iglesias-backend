package models

import (
	"gorm.io/gorm"
)

type Mensaje struct {
	gorm.Model
	Asunto        string              `json:"asunto"`
	Mensaje       string              `json:"mensaje"`
	Imagenes      string              `json:"imagenes_string"`
	ImagenesArray []string            `json:"imagenes" gorm:"-"`
	Estado        string              `gorm:"default:'NO_LEIDO'; gorm:type:enum('NO_LEIDO', 'LEIDO')"`
	AutorID       uint                `json:"id_usuario"`
	ParroquiaID   uint                `json:"id_etapa"`
	NoticiaID     *uint               `json:"id_noticia"`
	Autor         *Usuario            `json:"autor"`
	Noticia       *Publicacion        `json:"noticia"`
	Etapa         *Etapa              `json:"etapa"`
	Respuestas    []*RespuestaMensaje `json:"respuestas"`
}

func (Mensaje) TableName() string {
	return "mensaje"
}
