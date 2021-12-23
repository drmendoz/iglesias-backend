package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPadres(c *gin.Context) {
	idParroquia := c.GetInt("id_parroquia")
	padres := []*models.Padre{}
	err := models.Db.Order("nombre ASC").Where("parroquia_id = ?", idParroquia).Find(&padres).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener padres"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, padres, c, http.StatusOK)
}

func GetPadrePorId(c *gin.Context) {
	padre := &models.Padre{}
	id := c.Param("id")
	err := models.Db.First(padre, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Padre no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener padre"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, padre, c, http.StatusOK)
}

func CreatePadre(c *gin.Context) {
	padre := &models.Padre{}
	idParroquia := c.GetInt("id_parroquia")
	err := c.ShouldBindJSON(padre)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	padre.ParroquiaID = uint(idParroquia)
	tx := models.Db.Begin()
	err = tx.Create(padre).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear padre"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Padre creado correctamente", c, http.StatusCreated)
}

func UpdatePadre(c *gin.Context) {
	padre := &models.Padre{}

	err := c.ShouldBindJSON(padre)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Where("id = ?", id).Updates(padre).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar padre"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Padre actualizada correctamente", c, http.StatusOK)
}

func DeletePadre(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Padre{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar padre"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Padre eliminada exitosamente", c, http.StatusOK)
}
