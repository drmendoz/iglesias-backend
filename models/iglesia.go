package models

import "gorm.io/gorm"

type Iglesia struct {
	gorm.Model
	Nombre     string       `json:"nombre"`
	Parroquias []*Parroquia `json:"parroquias"`
}
