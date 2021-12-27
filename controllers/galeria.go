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

func GetGalerias(c *gin.Context) {
	idParroquia := c.GetInt("id_parroquia")
	imagens := []*models.Galeria{}
	err := models.Db.Where(&models.Galeria{ParroquiaID: uint(idParroquia)}).Find(&imagens).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener imagens"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, imagen := range imagens {
		if imagen.Imagen == "" {
			imagen.Imagen = utils.DefaultGaleria
		} else {
			imagen.Imagen = utils.SERVIMG + imagen.Imagen
		}

	}
	utils.CrearRespuesta(err, imagens, c, http.StatusOK)
}

func GetGaleriaPorId(c *gin.Context) {
	imagen := &models.Galeria{}
	id := c.Param("id")
	err := models.Db.First(imagen, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Galeria no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener galeria"), nil, c, http.StatusInternalServerError)
		return
	}
	if imagen.Imagen == "" {
		imagen.Imagen = utils.DefaultGaleria
	} else {
		imagen.Imagen = utils.SERVIMG + imagen.Imagen
	}

	utils.CrearRespuesta(nil, imagen, c, http.StatusOK)
}

func CreateGaleria(c *gin.Context) {
	imagen := &models.Galeria{}
	err := c.ShouldBindJSON(imagen)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	idParroquia := uint(c.GetInt("id_parroquia"))
	imagen.ParroquiaID = idParroquia
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(imagen).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear imagen"), nil, c, http.StatusInternalServerError)
		return
	}

	if imagen.Imagen == "" {
		imagen.Imagen = utils.DefaultGaleria
	} else {
		idUrb := fmt.Sprintf("%d", imagen.ID)
		imagen.Imagen, err = img.FromBase64ToImage(imagen.Imagen, "imagens/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(imagen.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear imagen "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Galeria{}).Where("id = ?", imagen.ID).Update("imagen", imagen.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear imagen "), nil, c, http.StatusInternalServerError)
			return
		}
		imagen.Imagen = utils.SERVIMG + imagen.Imagen
	}
	_ = tx.Commit()
	utils.CrearRespuesta(err, "Imagen creada correctamente", c, http.StatusCreated)

}

func UpdateGaleria(c *gin.Context) {
	imagen := &models.Galeria{}

	err := c.ShouldBindJSON(imagen)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen", "parroquia_id").Where("id = ?", id).Updates(imagen).Error
	if err != nil {
		_ = tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar imagen"), nil, c, http.StatusInternalServerError)
		return
	}
	if imagen.Imagen != "" {
		idUrb := fmt.Sprintf("%d", imagen.ID)
		imagen.Imagen, err = img.FromBase64ToImage(imagen.Imagen, "imagens/"+time.RFC3339+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear imagen "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Galeria{}).Where("id = ?", imagen.ID).Update("imagen", imagen.Imagen).Error
		if err != nil {
			_ = tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar imagen"), nil, c, http.StatusInternalServerError)
			return
		}
		imagen.Imagen = utils.SERVIMG + imagen.Imagen

	} else {
		imagen.Imagen = utils.DefaultGaleria
	}
	_ = tx.Commit()
	utils.CrearRespuesta(err, "Imagen actualizada correctamente", c, http.StatusOK)
}

func DeleteGaleria(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Galeria{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar imagen"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Imagen eliminada exitosamente", c, http.StatusOK)
}
