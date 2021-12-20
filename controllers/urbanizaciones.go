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
	"gorm.io/gorm/clause"
)

func GetUrbanizaciones(c *gin.Context) {
	urbs := []*models.Urbanizacion{}
	err := models.Db.Order("nombre ASC").Find(&urbs).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener urbanizaciones"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, urb := range urbs {
		if urb.Imagen == "" {
			urb.Imagen = utils.DefaultUrb
		} else {
			urb.Imagen = utils.SERVIMG + urb.Imagen
		}

	}
	utils.CrearRespuesta(err, urbs, c, http.StatusOK)
}

func GetUrbanizacionPorId(c *gin.Context) {

	urb := &models.Urbanizacion{}
	id := c.GetInt("id_urbanizacion")
	if id == 0 {

		id, _ = strconv.Atoi(c.Param("id"))
	}

	err := models.Db.Preload("Etapas").First(urb, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Urbanizacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener urbanizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	if urb.Imagen == "" {
		urb.Imagen = utils.DefaultUrb
	} else {
		urb.Imagen = utils.SERVIMG + urb.Imagen
	}
	for _, etp := range urb.Etapas {
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultEtapa
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}
	}
	utils.CrearRespuesta(nil, urb, c, http.StatusOK)
}

func CreateUrbanizacion(c *gin.Context) {
	urb := &models.Urbanizacion{}
	err := c.ShouldBindJSON(urb)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(urb).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear urbanizacion"), nil, c, http.StatusInternalServerError)
		return
	}

	if urb.Imagen == "" {
		urb.Imagen = utils.DefaultUrb
	} else {
		idUrb := fmt.Sprintf("%d", urb.ID)
		urb.Imagen, err = img.FromBase64ToImage(urb.Imagen, "urbanizaciones/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(urb.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear urbanizacion "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Urbanizacion{}).Where("id = ?", urb.ID).Update("imagen", urb.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear urbanizacion "), nil, c, http.StatusInternalServerError)
			return
		}
		urb.Imagen = utils.SERVIMG + urb.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Urbanizacion creada correctamente", c, http.StatusCreated)

}

func UpdateUrbanizacion(c *gin.Context) {
	urb := &models.Urbanizacion{}

	err := c.ShouldBindJSON(urb)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	idUrb := fmt.Sprintf("%d", urb.ID)
	if urb.Imagen != "" {
		urb.Imagen, err = img.FromBase64ToImage(urb.Imagen, "urbanizaciones/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(err, nil, c, http.StatusInternalServerError)
			return
		}
	}
	err = models.Db.Where("id = ?", id).Updates(urb).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar urbanizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, "Urbanizacion actualizada correctamente", c, http.StatusOK)
}

func DeleteUrbanizacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Urbanizacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar urbanizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Urbanizacion eliminada exitosamente", c, http.StatusOK)
}

func GetEtapasPorUrbanizacion(c *gin.Context) {
	urb := &models.Urbanizacion{}
	id := c.Param("id")
	err := models.Db.Preload(clause.Associations).First(urb, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Urbanizacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener urbanizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	if urb.Imagen == "" {
		urb.Imagen = utils.DefaultUrb
	} else {
		urb.Imagen = utils.SERVIMG + urb.Imagen
	}

	utils.CrearRespuesta(nil, urb, c, http.StatusOK)
}

func GetUrbanizacionesCount(c *gin.Context) {
	var urbs int64
	err := models.Db.Model(&models.Urbanizacion{}).Count(&urbs).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, urbs, c, http.StatusOK)
}
