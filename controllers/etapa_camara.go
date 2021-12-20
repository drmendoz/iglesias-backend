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

func GetEtapaCamaras(c *gin.Context) {
	camaras := []*models.EtapaCamara{}
	idEtapa := c.GetInt("id_etapa")
	err := models.Db.Where(&models.EtapaCamara{EtapaID: uint(idEtapa)}).Find(&camaras).Error

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener camaras"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, camara := range camaras {
		if camara.Imagen == "" {
			camara.Imagen = utils.DefaultCam
		} else {
			camara.Imagen = utils.SERVIMG + camara.Imagen
		}

	}
	utils.CrearRespuesta(err, camaras, c, http.StatusOK)
}

func GetEtapaCamaraPorId(c *gin.Context) {
	camara := &models.EtapaCamara{}
	id := c.Param("id")
	err := models.Db.First(camara, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Camara no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener camara"), nil, c, http.StatusInternalServerError)
		return
	}
	if camara.Imagen == "" {
		camara.Imagen = utils.DefaultCam
	} else {
		camara.Imagen = utils.SERVIMG + camara.Imagen
	}
	utils.CrearRespuesta(nil, camara, c, http.StatusOK)
}

func CreateEtapaCamara(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	camara := &models.EtapaCamara{}
	err := c.ShouldBindJSON(camara)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	camara.EtapaID = idEtapa
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(camara).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear camara"), nil, c, http.StatusInternalServerError)
		return
	}

	if camara.Imagen == "" {
		camara.Imagen = utils.DefaultCam
	} else {
		idUrb := fmt.Sprintf("%d", camara.ID)
		camara.Imagen, err = img.FromBase64ToImage(camara.Imagen, "camaras/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear camara "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.EtapaCamara{}).Where("id = ?", camara.ID).Update("imagen", camara.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear camara "), nil, c, http.StatusInternalServerError)
			return
		}
		camara.Imagen = utils.SERVIMG + camara.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Camara creada con exito", c, http.StatusCreated)

}

func UpdateEtapaCamara(c *gin.Context) {
	camara := &models.EtapaCamara{}

	err := c.ShouldBindJSON(camara)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(camara).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar camara"), nil, c, http.StatusInternalServerError)
		return
	}
	if camara.Imagen != "" {
		idUrb := fmt.Sprintf("%d", camara.ID)
		camara.Imagen, err = img.FromBase64ToImage(camara.Imagen, "camaras/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear camara "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.EtapaCamara{}).Where("id = ?", camara.ID).Update("imagen", camara.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al actualizar camara"), nil, c, http.StatusInternalServerError)
			return
		}
		camara.Imagen = utils.SERVIMG + camara.Imagen

	} else {
		camara.Imagen = utils.DefaultCam
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Camara actualizada correctamente", c, http.StatusOK)
}

func DeleteEtapaCamara(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.EtapaCamara{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar camara"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Camara eliminada exitosamente", c, http.StatusOK)
}
