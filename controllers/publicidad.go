package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPublicidads(c *gin.Context) {
	publicidads := []*models.Publicidad{}
	var err error

	idEtapa := c.GetInt("id_etapa")
	if idEtapa == 0 {
		idEtapa, _ = strconv.Atoi(c.Query("id_etapa"))
	}

	err = models.Db.Order("prioridad ASC").Where(&models.Publicidad{EtapaID: uint(idEtapa)}).Find(&publicidads).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener publicidads"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, publicidad := range publicidads {
		if publicidad.Imagen == "" {
			publicidad.Imagen = utils.DefaultPublicidad
		} else {
			publicidad.Imagen = utils.SERVIMG + publicidad.Imagen
		}

	}
	utils.CrearRespuesta(err, publicidads, c, http.StatusOK)
}

func GetPublicidadPorEtapaId(c *gin.Context) {
	publicidads := []*models.Publicidad{}
	err := models.Db.Order("prioridad ASC").Where("etapa_id = ?", c.Param("id")).Find(&publicidads).Error

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener publcidads"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, publicidad := range publicidads {
		if publicidad.Imagen == "" {
			publicidad.Imagen = utils.DefaultPublicidad
		} else {
			publicidad.Imagen = utils.SERVIMG + publicidad.Imagen
		}
	}
	utils.CrearRespuesta(err, publicidads, c, http.StatusOK)

}

func GetPublicidadPorId(c *gin.Context) {
	publicidad := &models.Publicidad{}
	id := c.Param("id")
	err := models.Db.First(publicidad, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Publicidad no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener publicidad"), nil, c, http.StatusInternalServerError)
		return
	}
	if publicidad.Imagen == "" {
		publicidad.Imagen = utils.DefaultPublicidad
	} else {
		publicidad.Imagen = utils.SERVIMG + publicidad.Imagen
	}

	utils.CrearRespuesta(nil, publicidad, c, http.StatusOK)
}

func CreatePublicidad(c *gin.Context) {
	publicidad := &models.Publicidad{}
	err := c.ShouldBindJSON(publicidad)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(publicidad).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear publicidad"), nil, c, http.StatusInternalServerError)
		return
	}

	if publicidad.Imagen == "" {
		publicidad.Imagen = utils.DefaultPublicidad
	} else {
		idUrb := fmt.Sprintf("%d", publicidad.ID)
		publicidad.Imagen, err = img.FromBase64ToImage(publicidad.Imagen, "publicidads/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(publicidad.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear publicidad "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Publicidad{}).Where("id = ?", publicidad.ID).Update("imagen", publicidad.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear publicidad "), nil, c, http.StatusInternalServerError)
			return
		}
		publicidad.Imagen = utils.SERVIMG + publicidad.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Publicidad creada con exito", c, http.StatusCreated)

}

func UpdatePublicidad(c *gin.Context) {
	publicidad := &models.Publicidad{}

	err := c.ShouldBindJSON(publicidad)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(publicidad).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar publicidad"), nil, c, http.StatusInternalServerError)
		return
	}
	println(img.IsBase64(publicidad.Imagen))
	if img.IsBase64(publicidad.Imagen) {
		idUrb := fmt.Sprintf("%d", publicidad.ID)
		publicidad.Imagen, err = img.FromBase64ToImage(publicidad.Imagen, "publicidads/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear publicidad "), nil, c, http.StatusInternalServerError)
			return
		}
		err = tx.Model(&models.Publicidad{}).Where("id = ?", id).Update("imagen", publicidad.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar publicidad"), nil, c, http.StatusInternalServerError)
			return
		}
		publicidad.Imagen = utils.SERVIMG + publicidad.Imagen

	} else {
		publicidad.Imagen = utils.DefaultPublicidad
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Publicidad actualizada correctamente", c, http.StatusOK)
}

func DeletePublicidad(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Publicidad{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar publicidad"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Publicidad eliminada exitosamente", c, http.StatusOK)
}
