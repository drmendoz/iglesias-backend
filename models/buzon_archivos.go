package models

import "gorm.io/gorm"

type BuzonArchivo struct {
	gorm.Model
	Url      string `json:"url"`
	MimeType string `json:"mime"`
	BuzonID  uint   `json:"id_buzon"`
	Buzon    *Buzon `json:"buzon"`
}

func (BuzonArchivo) TableName() string {
	return "buzon_archivo"
}
