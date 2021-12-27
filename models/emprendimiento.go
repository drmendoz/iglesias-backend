package models

import (
	"gorm.io/gorm"
)

type Emprendimiento struct {
	gorm.Model
	Titulo                 string                  `json:"titulo"`
	Descripcion            string                  `json:"descripcion"`
	Precio                 float64                 `json:"-"`
	Estado                 string                  `json:"estado" gorm:"default:'DES';type:enum('VIG','DES')"`
	Imagen                 string                  `json:"imagen,omitempty" gorm:"-"`
	TelefonoContacto       string                  `json:"telefono_contacto"`
	PrecioLabel            string                  `json:"precio"`
	NombreUsuario          string                  `json:"nombre_usuario" gorm:"-"`
	ImagenUsuario          string                  `json:"imagen_usuario" gorm:"-"`
	TelefonoUsuario        string                  `json:"celular_usuario" gorm:"-"`
	Ciudad                 string                  `json:"ciudad"`
	EmprendimientoImagenes []*EmprendimientoImagen `json:"-" `
	Imagenes               []string                `json:"imagenes" gorm:"-"`
	TokenTarjeta           string                  `json:"token_tarjeta" gorm:"-"`
	//FechaPublicacion       time.Time               `json:"fecha_publicacion"`
	//	FechaVencimiento       time.Time               `json:"fecha_vencimiento"`
	FielID            uint             `json:"id_fiel"`
	ParroquiaID       uint             `json:"id_parroquia"`
	Transaccion       *Transaccion     `json:"transaccion" gorm:"polymorphic:TipoPago"`
	CategoriaMarketID uint             `json:"id_categoria"`
	CategoriaMarket   *CategoriaMarket `json:"categoria,omitempty"`
	Parroquia         *Parroquia       `json:"parroquia,omitempty"`
	Fiel              *Fiel            `json:"fiel,omitempty"`
}

func (Emprendimiento) TableName() string {
	return "emprendimiento"
}
