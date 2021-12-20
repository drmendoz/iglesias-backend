package models

import (
	"time"

	"gorm.io/gorm"
)

type Emprendimiento struct {
	gorm.Model
	Titulo                 string                  `json:"titulo"`
	Descripcion            string                  `json:"descripcion"`
	Precio                 float64                 `json:"-"`
	Estado                 string                  `json:"estado" gorm:"default:'VIG';type:enum('VIG','DES')"`
	Imagen                 string                  `json:"imagen,omitempty" gorm:"-"`
	TelefonoContacto       string                  `json:"telefono_contacto"`
	PrecioLabel            string                  `json:"precio"`
	NombreUsuario          string                  `json:"nombre_usuario" gorm:"-"`
	ImagenUsuario          string                  `json:"imagen_usuario" gorm:"-"`
	TelefonoUsuario        string                  `json:"celular_usuario" gorm:"-"`
	Ciudad                 string                  `json:"ciudad"`
	EmprendimientoImagenes []*EmprendimientoImagen `json:"-" `
	FechaPublicacion       time.Time               `json:"fecha_publicacion"`
	FechaVencimiento       time.Time               `json:"fecha_vencimiento"`
	Premium                bool                    `json:"premium" `
	ResidenteID            uint                    `json:"id_residente"`
	EtapaID                uint                    `json:"id_etapa"`
	CategoriaMarketID      uint                    `json:"id_categoria"`
	CategoriaMarket        *CategoriaMarket        `json:"categoria,omitempty"`
	Residente              *Residente              `json:"residente,omitempty"`
	Imagenes               []string                `json:"imagenes" gorm:"-"`
}

func (Emprendimiento) TableName() string {
	return "emprendimiento"
}
