package models

import "gorm.io/gorm"

type Buzon struct {
	gorm.Model
	Titulo           string               `json:"titulo" `
	Descripcion      string               `json:"descripcion" `
	Leido            bool                 `json:"leido" gorm:"-"`
	PublicadorID     uint                 `json:"publicador_id"`
	Publico          bool                 `json:"publico" gorm:"default:false"`
	Adjuntos         bool                 `json:"adjuntos" gorm:"-"`
	IsAdmin          bool                 `json:"is_admin" gorm:"default:false"`
	CasaID           *uint                `json:"id_casa"`
	EsRespuesta      bool                 `json:"es_respuesta" gorm:"default:false"`
	Casa             *Casa                `json:"casa"`
	Publicador       *Usuario             `json:"publicador,omitempty" gorm:"foreignKey:PublicadorID"`
	Destinatarios    []*BuzonDestinatario `json:"destinatarios"`
	Archivos         []*BuzonArchivo      `json:"archivos"`
	UltimoMensaje    *Buzon               `json:"ultimo_mensaje" gorm:"-"`
	BuzonRemitenteID *uint                `json:"buzon_remitente_id"`
	BuzonRemitente   *Buzon               `json:"buzon_remitente"`
	Mensajes         []*Buzon             `json:"mensajes" gorm:"foreignKey:BuzonRemitenteID"`
	EtapaID          uint                 `json:"id_etapa"`
	Etapa            *Etapa               `json:"etapa"`
}

func (Buzon) TableName() string {
	return "buzon"
}
