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

func GetImagenGalerias(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	imagens := []*models.ImagenGaleria{}
	var err error
	if idParroquia != 0 {
		err = models.Db.Where("etapa_id = ?", idParroquia).Order("created_at desc").Find(&imagens).Error
	} else {
		err = models.Db.Find(&imagens).Error
	}
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

func GetImagenGaleriaPorId(c *gin.Context) {
	imagen := &models.ImagenGaleria{}
	id := c.Param("id")
	err := models.Db.First(imagen, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Imagen galeria no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener imagen"), nil, c, http.StatusInternalServerError)
		return
	}
	if imagen.Imagen == "" {
		imagen.Imagen = utils.DefaultGaleria
	} else {
		imagen.Imagen = utils.SERVIMG + imagen.Imagen
	}

	utils.CrearRespuesta(nil, imagen, c, http.StatusOK)
}

func CreateImagenGaleria(c *gin.Context) {
	imagen := &models.ImagenGaleria{}
	err := c.ShouldBindJSON(imagen)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	idParroquia := uint(c.GetInt("id_etapa"))
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
		err = tx.Model(&models.ImagenGaleria{}).Where("id = ?", imagen.ID).Update("imagen", imagen.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear imagen "), nil, c, http.StatusInternalServerError)
			return
		}
		imagen.Imagen = utils.SERVIMG + imagen.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Imagen galeria creada correctamente", c, http.StatusCreated)

}

func UpdateImagenGaleria(c *gin.Context) {
	imagen := &models.ImagenGaleria{}

	err := c.ShouldBindJSON(imagen)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen", "etapa_id").Where("id = ?", id).Updates(imagen).Error
	if err != nil {
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
		err = tx.Model(&models.ImagenGaleria{}).Where("id = ?", imagen.ID).Update("imagen", imagen.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar imagen"), nil, c, http.StatusInternalServerError)
			return
		}
		imagen.Imagen = utils.SERVIMG + imagen.Imagen

	} else {
		imagen.Imagen = utils.DefaultGaleria
	}
	utils.CrearRespuesta(err, "Imagen galeria actualizada correctamente", c, http.StatusOK)
}

func DeleteImagenGaleria(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.ImagenGaleria{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar imagen"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Imagen galeria eliminada exitosamente", c, http.StatusOK)
}
