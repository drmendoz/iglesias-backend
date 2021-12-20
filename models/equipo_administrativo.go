package models

type EquipoAdministrativo struct {
	Nombres     string `json:"nombres"`
	Cargo       string `json:"cargo"`
	Imagen      string `json:"imagen"`
	Cedula      string `json:"cedula"`
	Correo      string `json:"correo"`
	Tipo        string `json:"tipo" gorm:"default:'ADMINISTRATIVO';type:enum('ADMINISTRATIVO','OPERATIVO')"`
	ParroquiaID uint   `json:"id_etapa"`
	Etapa       Etapa  `json:"etapa"`
}
