package models

import "gorm.io/gorm"

type Padre struct {
	gorm.Model
	ParroquiaID uint       `json:"id_parroquia"`
	Nombre      uint       `json:"nombre"`
	Parroquia   *Parroquia `json:"parroquia"`
}

func (Padre) TableName() string {
	return "padre"
}
