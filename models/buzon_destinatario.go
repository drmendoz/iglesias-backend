package models

import "gorm.io/gorm"

type BuzonDestinatario struct {
	gorm.Model
	CasaID  uint   `json:"id_casa"`
	Casa    *Casa  `json:"casa"`
	Rol     string `json:"rol" gorm:"type:enum('admin-etapa','residente')"`
	Leido   bool   `json:"leido" gorm:"default:false"`
	BuzonID uint   `json:"id_buzon"`
	Buzon   *Buzon `json:"buzon"`
}

func (BuzonDestinatario) TableName() string {
	return "buzon_destinatario"
}
