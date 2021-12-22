package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetIntenciones(c *gin.Context) {
	etps := []*models.Intencion{}
	err := models.Db.Order("created_at ASC").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener intencions"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetIntencionPorId(c *gin.Context) {
	etp := &models.Intencion{}
	id := c.Param("id")
	err := models.Db.First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Intención no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener intención"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateIntencion(c *gin.Context) {
	etp := &models.Intencion{}
	idParroquia := c.GetInt("id_parroquia")
	idFiel := c.GetInt("id_fiel")

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = uint(idParroquia)
	etp.FielID = uint(idFiel)

	tx := models.Db.Begin()
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear intención"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Intención creada correctamente", c, http.StatusCreated)

}

func UpdateIntencion(c *gin.Context) {
	etp := &models.Intencion{}
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
		utils.CrearRespuesta(errors.New("Error al actualizar intención"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Intención actualizada correctamente", c, http.StatusOK)
}

func DeleteIntencion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Intencion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar intención"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Intención eliminada exitosamente", c, http.StatusOK)
}
