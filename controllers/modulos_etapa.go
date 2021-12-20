package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetModuloEtapas(c *gin.Context) {
	modulos := []*models.ModuloEtapa{}
	idEtapa := c.GetInt("id_etapa")
	err := models.Db.Where(&models.ModuloEtapa{EtapaID: uint(idEtapa)}).Find(&modulos).Error

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener modulos"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, modulo := range modulos {
		if modulo.Imagen == "" {
			modulo.Imagen = utils.DefaultCam
		} else {
			modulo.Imagen = utils.SERVIMG + modulo.Imagen
		}

	}
	utils.CrearRespuesta(err, modulos, c, http.StatusOK)
}

func GetModuloEtapaPorId(c *gin.Context) {
	modulo := &models.ModuloEtapa{}
	id := c.Param("id")
	err := models.Db.First(modulo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Modulo no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener modulo"), nil, c, http.StatusInternalServerError)
		return
	}
	if modulo.Imagen == "" {
		modulo.Imagen = utils.DefaultCam
	} else {
		modulo.Imagen = utils.SERVIMG + modulo.Imagen
	}
	utils.CrearRespuesta(nil, modulo, c, http.StatusOK)
}

func CreateModuloEtapa(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	modulo := &models.ModuloEtapa{}
	err := c.ShouldBindJSON(modulo)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	modulo.EtapaID = idEtapa
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(modulo).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear modulo"), nil, c, http.StatusInternalServerError)
		return
	}

	if modulo.Imagen == "" {
		modulo.Imagen = utils.DefaultCam
	} else {
		idUrb := fmt.Sprintf("%d", modulo.ID)
		modulo.Imagen, err = img.FromBase64ToImage(modulo.Imagen, "modulos/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear modulo "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.ModuloEtapa{}).Where("id = ?", modulo.ID).Update("imagen", modulo.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear modulo "), nil, c, http.StatusInternalServerError)
			return
		}
		modulo.Imagen = utils.SERVIMG + modulo.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Modulo creada con exito", c, http.StatusCreated)

}

func UpdateModuloEtapa(c *gin.Context) {
	modulo := &models.ModuloEtapa{}

	err := c.ShouldBindJSON(modulo)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(modulo).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar modulo"), nil, c, http.StatusInternalServerError)
		return
	}
	if modulo.Imagen != "" {
		idUrb := fmt.Sprintf("%d", modulo.ID)
		modulo.Imagen, err = img.FromBase64ToImage(modulo.Imagen, "modulos/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear modulo "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.ModuloEtapa{}).Where("id = ?", modulo.ID).Update("imagen", modulo.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al actualizar modulo"), nil, c, http.StatusInternalServerError)
			return
		}
		modulo.Imagen = utils.SERVIMG + modulo.Imagen

	} else {
		modulo.Imagen = utils.DefaultCam
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Modulo actualizada correctamente", c, http.StatusOK)
}

func DeleteModuloEtapa(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.ModuloEtapa{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar modulo"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Modulo eliminada exitosamente", c, http.StatusOK)
}
