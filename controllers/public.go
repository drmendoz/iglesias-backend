package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
)

type Archivo struct {
	Archivo string `json:"archivo"`
}

func CrearArchivo(c *gin.Context) {
	archivo := &Archivo{}
	err := c.ShouldBindJSON(archivo)
	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros de Request Invalido"), nil, c, http.StatusBadRequest)
		return
	}
	imagen, err := img.FromBase64ToImage(archivo.Archivo, "emprendimiento/"+time.Now().Format(time.RFC3339), false)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear imagen"), nil, c, http.StatusInternalServerError)
		return
	}
	imagen = utils.SERVIMG + imagen
	utils.CrearRespuesta(nil, imagen, c, http.StatusCreated)
}

type Archivos struct {
	Archivos []string `json:"archivos"`
}

func CrearArchivos(c *gin.Context) {
	archivos := &Archivos{}
	err := c.ShouldBindJSON(archivos)
	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros de Request Invalido"), nil, c, http.StatusBadRequest)
		return
	}
	urls := []string{}
	cont := 0
	for _, arc := range archivos.Archivos {
		ct := fmt.Sprintf("%d", cont)
		imagen, err := img.FromBase64ToImage(arc, "emprendimiento/"+time.Now().Format(time.RFC3339)+ct, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear imagen"), nil, c, http.StatusInternalServerError)
			return
		}
		imagen = utils.SERVIMG + imagen
		urls = append(urls, imagen)

		cont++
	}

	utils.CrearRespuesta(nil, urls, c, http.StatusCreated)
}
