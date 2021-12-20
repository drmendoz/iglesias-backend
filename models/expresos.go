package models

type Empleados struct {
	Nombre string `json:"nombre"`
	Cedula string `json:"cedula"`
	CasaID uint   `json:"id_cas"`
	Casa   *Casa  `json:"casa"`
}
