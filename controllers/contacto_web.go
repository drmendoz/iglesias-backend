package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/mail"
	"github.com/gin-gonic/gin"
)

func SolicitudContacto(c *gin.Context) {
	contacto := mail.Contacto{}
	err := c.ShouldBindJSON(&contacto)
	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	err = mail.EnviarCorreoContactoWeb(contacto)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al enviar solicitud de contacto"), nil, c, http.StatusInsufficientStorage)
		return
	}
	utils.CrearRespuesta(nil, "Formulario enviado exitosamente", c, http.StatusOK)
}
