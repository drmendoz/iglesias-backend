package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRespuestaVotacions(c *gin.Context) {
	respuestas := []*models.RespuestaVotacion{}
	err := models.Db.Find(&respuestas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener respuestas"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, respuestas, c, http.StatusOK)
}

func GetRespuestaVotacionPorId(c *gin.Context) {
	respuesta := &models.RespuestaVotacion{}
	id := c.Param("id")
	err := models.Db.First(respuesta, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Respuesta no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener respuesta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, respuesta, c, http.StatusOK)
}

func CreateRespuestaVotacion(c *gin.Context) {
	respuesta := &models.RespuestaVotacion{}
	err := c.ShouldBindJSON(respuesta)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Create(respuesta).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear respuesta"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Respuesta votacion creada correctamente", c, http.StatusCreated)

}

func UpdateRespuestaVotacion(c *gin.Context) {
	respuesta := &models.RespuestaVotacion{}

	err := c.ShouldBindJSON(respuesta)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(respuesta).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar respuesta"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Respuesta votacion actualizada correctamente", c, http.StatusOK)
}

func DeleteRespuestaVotacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.RespuestaVotacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar respuesta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Respuesta votacion eliminada exitosamente", c, http.StatusOK)
}
