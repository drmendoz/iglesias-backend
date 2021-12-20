package models

import "gorm.io/gorm"

type Casa struct {
	gorm.Model
	Manzana        string      `json:"manzana"`
	Villa          string      `json:"villa"`
	Direccion      string      `json:"direccion"`
	Piso           string      `json:"piso"`
	Condominio     string      `json:"departamento"`
	TipoCasa       string      `json:"tipo_casa" gorm:"default:'CASA';type:enum('CASA','DEP')"`
	ParroquiaID    uint        `json:"id_etapa"`
	Familia        string      `json:"familia"`
	Imagen         string      `json:"imagen"`
	Fijo           string      `json:"fijo"`
	Celular        string      `json:"celular"`
	DebeAlicuotas  bool        `json:"debe_alicuotas" gorm:"-"`
	ValorAlicuotas float64     `json:"valor_alicuotas" gorm:"-"`
	Etapa          *Etapa      `json:"etapa,omitempty"`
	Fiels          []*Fiel     `json:"residentes,omitempty"`
	Alicuotas      []*Alicuota `json:"alicuotas,omitempty"`
}

func (Casa) TableName() string {
	return "casa"
}
