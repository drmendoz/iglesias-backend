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

func GetAdministrativos(c *gin.Context) {
	administrativos := []*models.Administrativo{}
	var err error
	idParroquia := c.GetInt("id_etapa")
	if idParroquia != 0 {
		err = models.Db.Where("etapa_id = ?", idParroquia).Find(&administrativos).Error
	} else {
		err = models.Db.Find(&administrativos).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administrativos"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, administrativo := range administrativos {
		if administrativo.Imagen == "" {
			administrativo.Imagen = utils.DefaultAdministrativo
		} else {
			administrativo.Imagen = utils.SERVIMG + administrativo.Imagen
		}
	}
	utils.CrearRespuesta(err, administrativos, c, http.StatusOK)
}

func GetAdministrativoPorId(c *gin.Context) {
	administrativo := &models.Administrativo{}
	id := c.Param("id")
	err := models.Db.First(administrativo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Administrativo no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administrativo"), nil, c, http.StatusInternalServerError)
		return
	}
	if administrativo.Imagen == "" {
		administrativo.Imagen = utils.DefaultAdministrativo
	} else {
		administrativo.Imagen = utils.SERVIMG + administrativo.Imagen
	}

	utils.CrearRespuesta(nil, administrativo, c, http.StatusOK)
}

func CreateAdministrativo(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_etapa"))
	administrativo := &models.Administrativo{}
	err := c.ShouldBindJSON(administrativo)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	administrativo.ParroquiaID = idParroquia
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(administrativo).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear administrativo"), nil, c, http.StatusInternalServerError)
		return
	}

	if administrativo.Imagen == "" {
		administrativo.Imagen = utils.DefaultAdministrativo
	} else {
		idUrb := fmt.Sprintf("%d", administrativo.ID)
		administrativo.Imagen, err = img.FromBase64ToImage(administrativo.Imagen, "administrativos/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(administrativo.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear administrativo "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Administrativo{}).Where("id = ?", administrativo.ID).Update("imagen", administrativo.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear administrativo "), nil, c, http.StatusInternalServerError)
			return
		}
		administrativo.Imagen = utils.SERVIMG + administrativo.Imagen
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Administrativo creada exitosamente", c, http.StatusCreated)

}

func UpdateAdministrativo(c *gin.Context) {
	administrativo := &models.Administrativo{}

	err := c.ShouldBindJSON(administrativo)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("etapa_id", "imagen").Where("id = ?", id).Updates(administrativo).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar administrativo"), nil, c, http.StatusInternalServerError)
		return
	}
	if img.IsBase64(administrativo.Imagen) {
		idUrb := fmt.Sprintf("%d", administrativo.ID)
		administrativo.Imagen, err = img.FromBase64ToImage(administrativo.Imagen, "administrativos/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear administrativo "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Administrativo{}).Where("id = ?", id).Update("imagen", administrativo.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al actualizar administrativo"), nil, c, http.StatusInternalServerError)
			return
		}
		administrativo.Imagen = utils.SERVIMG + administrativo.Imagen

	} else {
		administrativo.Imagen = utils.DefaultAdministrativo
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Administrativo actualizada correctamente", c, http.StatusOK)
}

func DeleteAdministrativo(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Administrativo{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar administrativo"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Administrativo eliminada exitosamente", c, http.StatusOK)
}
