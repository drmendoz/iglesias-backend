package utils

import (
	"github.com/gin-gonic/gin"
)

//Respuesta formato de envio de respuestas
type Respuesta struct {
	HasError  bool        `json:"error"`
	Contenido interface{} `json:"respuesta"`
}

//CrearRespuesta formato para enviar respuestas
func CrearRespuesta(Error error, Contenido interface{}, c *gin.Context, codigoError int) {
	respuesta := &Respuesta{}
	if Error != nil {
		respuesta.HasError = true
		respuesta.Contenido = Error.Error()

	} else {
		respuesta.HasError = false
		respuesta.Contenido = Contenido
	}
	c.Header("Server", "REST Api Urbanizaciones")
	c.JSON(codigoError, respuesta)
}
