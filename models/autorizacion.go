package models

import "gorm.io/gorm"

type Autorizacion struct {
	gorm.Model
	Nombre       string      `json:"nombres"`
	Apellidos    string      `json:"apellidos"`
	Cedula       string      `json:"cedula"`
	Imagen       string      `json:"imagen"`
	Imagen2      string      `json:"imagen2"`
	Motivo       string      `json:"asunto"`
	Manzana      string      `json:"mz"`
	Villa        string      `json:"villa"`
	Tipo         string      `json:"tipo" gorm:"type:enum('TEMPORAL','FIJA')"`
	TipoUsuario  string      `json:"tipo_usuario" gorm:"default:'VISITA'; gorm:type:enum('RESIDENTE','EMPLEADO','VISITA', 'EXPRESO', 'FAMILIAR')"`
	Estado       string      `json:"estado" gorm:"default:'PENDIENTE'; gorm:type:enum('ACTIVA', 'ANULADA', 'VALIDADA', 'PENDIENTE')"` //`json:"estado" gorm:"default:'PENDIENTE';type:enum('PENDIENTE','ACTIVA','ANULADA', 'USADA')"`
	Telefono     string      `json:"telefono"`
	Pin          string      `json:"pin"`
	Pdf          string      `json:"documento"`
	CasaID       uint        `json:"id_casa"`
	PublicadorID uint        `json:"id_guardia"`
	UsuarioID    uint        `json:"id_usuario"`
	ParroquiaID  uint        `json:"id_etapa"`
	AutorizadoID *uint       `json:"id_autorizado"`
	Casa         *Casa       `json:"casa,omitempty"`
	Etapa        *Etapa      `json:"etapa,omitempty"`
	Publicador   *Usuario    `json:"guardia,omitempty" gorm:"foreignKey:PublicadorID"`
	Usuario      *Usuario    `json:"contestador"`
	Autorizado   *Autorizado `json:"autorizado"`
	Correo       string      `json:"correo"`
	// Username     string   `json:"usuario"` // Usuario.Usuario
	// Nickname     string   `json:"usuario"`
}

func (Autorizacion) TableName() string {
	return "autorizacion"
}
