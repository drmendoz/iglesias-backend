package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMisas(c *gin.Context) {
	etps := []*models.Misa{}
	err := models.Db.Order("created_at ASC").Preload("Padre").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener misas"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetMisaPorId(c *gin.Context) {
	etp := &models.Misa{}
	id := c.Param("id")
	err := models.Db.First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Misa no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener misa"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateMisa(c *gin.Context) {
	etp := &models.Misa{}
	idParroquia := c.GetInt("id_parroquia")

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = uint(idParroquia)

	tx := models.Db.Begin()
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear misa"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Misa creada correctamente", c, http.StatusCreated)

}

func UpdateMisa(c *gin.Context) {
	etp := &models.Misa{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar misa"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Misa actualizada correctamente", c, http.StatusOK)
}

func DeleteMisa(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Misa{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar misa"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Misa eliminada exitosamente", c, http.StatusOK)
}
