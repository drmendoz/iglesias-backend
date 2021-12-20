package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOpcionVotacions(c *gin.Context) {
	opcions := []*models.OpcionVotacion{}
	err := models.Db.Find(&opcions).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener opcions"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, opcions, c, http.StatusOK)
}

func GetOpcionVotacionPorId(c *gin.Context) {
	opcion := &models.OpcionVotacion{}
	id := c.Param("id")
	err := models.Db.First(opcion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Opcion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener opcion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, opcion, c, http.StatusOK)
}

func CreateOpcionVotacion(c *gin.Context) {
	opcion := &models.OpcionVotacion{}
	err := c.ShouldBindJSON(opcion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Create(opcion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear opcion"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "OpcionVotacion creada correctamente", c, http.StatusCreated)

}

func UpdateOpcionVotacion(c *gin.Context) {
	opcion := &models.OpcionVotacion{}

	err := c.ShouldBindJSON(opcion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(opcion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar opcion"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "OpcionVotacion actualizada correctamente", c, http.StatusOK)
}

func DeleteOpcionVotacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.OpcionVotacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar opcion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "OpcionVotacion eliminada exitosamente", c, http.StatusOK)
}
