package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDeliverys(c *gin.Context) {
	respuestas := []*models.Delivery{}
	err := models.Db.Find(&respuestas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener respuestas"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, respuestas, c, http.StatusOK)
}

func GetDeliveryPorId(c *gin.Context) {
	respuesta := &models.Delivery{}
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

func CreateDelivery(c *gin.Context) {
	respuesta := &models.Delivery{}
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

	utils.CrearRespuesta(err, "Delivery creado correctamente", c, http.StatusCreated)

}

func UpdateDelivery(c *gin.Context) {
	respuesta := &models.Delivery{}

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

	utils.CrearRespuesta(err, "Delivery actualizado correctamente", c, http.StatusOK)
}

func DeleteDelivery(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Delivery{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar respuesta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Delivery eliminado exitosamente", c, http.StatusOK)
}
