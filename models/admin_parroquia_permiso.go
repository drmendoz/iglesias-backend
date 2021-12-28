package models

import "gorm.io/gorm"

type AdminParroquiaPermiso struct {
	gorm.Model
	AdminParroquiaID uint            `json:"id_admin_etapa"`
	AdminParroquia   *AdminParroquia `json:"admin_etapa"`
	Usuario          bool            `json:"usuario" gorm:"default:false"`
	Horario          bool            `json:"horario" gorm:"default:false"`
	Actividad        bool            `json:"actividad" gorm:"default:false"`
	Emprendimiento   bool            `json:"emprendimiento" gorm:"default:false"`
	Intencion        bool            `json:"intencion" gorm:"default:false"`
	Musica           bool            `json:"musica" gorm:"default:false"`
	Ayudemos         bool            `json:"ayudemos" gorm:"default:false"`
	Misa             bool            `json:"misa" gorm:"default:false"`
	Curso            bool            `json:"curso" gorm:"default:false"`
	Matrimonio       bool            `json:"matrimonio" gorm:"default:false"`
	Publicacion      bool            `json:"publicacion" gorm:"default:false"`
	Galeria          bool            `json:"galeria" gorm:"default:false"`
	Recaudacion      bool            `json:"recaudacion" gorm:"default:false"`
}

func (AdminParroquiaPermiso) TableName() string {
	return "admin_parroquia_permiso"
}
