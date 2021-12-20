package models

import "gorm.io/gorm"

type AdminEtapaPermiso struct {
	gorm.Model
	AdminEtapaID   uint        `json:"id_admin_etapa"`
	AdminEtapa     *AdminEtapa `json:"admin_etapa"`
	Alicuota       bool        `json:"alicuota" gorm:"default:false"`
	AreaSocial     bool        `json:"area_social" gorm:"default:false"`
	Emprendimiento bool        `json:"emprendimiento" gorm:"default:false"`
	Casa           bool        `json:"casa" gorm:"default:false"`
	Galeria        bool        `json:"galeria" gorm:"default:false"`
	Usuario        bool        `json:"usuario" gorm:"default:false"`
	Seguridad      bool        `json:"seguridad" gorm:"default:false"`
	Ingreso        bool        `json:"ingreso" gorm:"default:false"`
	Voto           bool        `json:"voto" gorm:"default:false"`
	Directiva      bool        `json:"directiva" gorm:"default:false"`
	Camara         bool        `json:"camara" gorm:"default:false"`
	Horario        bool        `json:"horario" gorm:"default:false"`
	Reserva        bool        `json:"reserva" gorm:"default:false"`
	ExpresoEscolar bool        `json:"expreso_escolar" gorm:"default:false"`
	Buzon          bool        `json:"buzon" gorm:"default:false"`
}

func (AdminEtapaPermiso) TableName() string {
	return "admin_etapa_permiso"
}
