package models

import (
	"gorm.io/gorm"
)

type RespuestaMensaje struct {
	gorm.Model
	Asunto        string       `json:"asunto"`
	Respuesta     string       `json:"mensaje"`
	Imagenes      string       `json:"imagenes_string"`
	ImagenesArray []string     `json:"imagenes" gorm:"-"`
	Estado        string       `gorm:"default:'NO_LEIDO'; gorm:type:enum('NO_LEIDO', 'LEIDO')"`
	AutorID       uint         `json:"id_usuario"`
	MensajeID     uint         `json:"id_mensaje"`
	NoticiaID     *uint        `json:"id_noticia"`
	Autor         *Usuario     `json:"autor"`
	Noticia       *Publicacion `json:"noticia"`
	Mensaje       *Mensaje     `json:"conversacion"`
}

func (RespuestaMensaje) TableName() string {
	return "respuesta_mensaje"
}
