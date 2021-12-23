package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMusicas(c *gin.Context) {
	etps := []*models.Musica{}
	idParroquia := c.GetInt("id_parroquia")
	err := models.Db.Where(&models.Musica{ParroquiaID: uint(idParroquia)}).Order("created_at asc").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener areas sociales"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, etps, c, http.StatusOK)
}

func GetMusicaPorID(c *gin.Context) {
	etp := &models.Musica{}
	id := c.Param("id")
	err := models.Db.Preload("Aportaciones").Preload("Aportaciones.Fiel").Preload("Aportaciones.Transaccion").First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Doncacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener donacion"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateMusica(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_parroquia"))
	etp := &models.Musica{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = idParroquia

	tx := models.Db.Begin()
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear donacion"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Musica creada correctamente", c, http.StatusCreated)

}

func UpdateMusica(c *gin.Context) {
	etp := &models.Musica{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	id := c.Param("id")
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Musica actualizada correctamente", c, http.StatusOK)
}

func DeleteMusica(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Musica{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Musica eliminada exitosamente", c, http.StatusOK)
}
