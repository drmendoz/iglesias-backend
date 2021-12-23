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

func GetCategoriaDonacions(c *gin.Context) {
	categorias := []*models.CategoriaDonacion{}
	err := models.Db.Find(&categorias).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener categorias"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, categoria := range categorias {
		if categoria.Imagen == "" {
			categoria.Imagen = utils.DefaultCategoria
		} else {
			categoria.Imagen = utils.SERVIMG + categoria.Imagen
		}

	}
	utils.CrearRespuesta(err, categorias, c, http.StatusOK)
}

func GetCategoriaDonacionPorId(c *gin.Context) {
	categoria := &models.CategoriaDonacion{}
	id := c.Param("id")
	err := models.Db.First(categoria, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Categoria no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener categoria"), nil, c, http.StatusInternalServerError)
		return
	}
	if categoria.Imagen == "" {
		categoria.Imagen = utils.DefaultCategoria
	} else {
		categoria.Imagen = utils.SERVIMG + categoria.Imagen
	}

	utils.CrearRespuesta(nil, categoria, c, http.StatusOK)
}

func CreateCategoriaDonacion(c *gin.Context) {
	categoria := &models.CategoriaDonacion{}
	err := c.ShouldBindJSON(categoria)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(categoria).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear categoria"), nil, c, http.StatusInternalServerError)
		return
	}

	if categoria.Imagen == "" {
		categoria.Imagen = utils.SERVIMG + "default_user.png"
	} else if img.IsBase64(categoria.Imagen) {
		idUrb := fmt.Sprintf("%d", categoria.ID)
		categoria.Imagen, err = img.FromBase64ToImage(categoria.Imagen, "categorias/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(categoria.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear categoria "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.CategoriaDonacion{}).Where("id = ?", categoria.ID).Update("imagen", categoria.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear categoria "), nil, c, http.StatusInternalServerError)
			return
		}
		categoria.Imagen = utils.SERVIMG + categoria.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Categoria creada correctamente", c, http.StatusCreated)

}

func UpdateCategoriaDonacion(c *gin.Context) {
	categoria := &models.CategoriaDonacion{}

	err := c.ShouldBindJSON(categoria)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(categoria).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar categoria"), nil, c, http.StatusInternalServerError)
		return
	}
	if img.IsBase64(categoria.Imagen) {
		idUrb := fmt.Sprintf("%d", categoria.ID)
		categoria.Imagen, err = img.FromBase64ToImage(categoria.Imagen, "categorias/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			return
		}
		err = tx.Model(&models.CategoriaDonacion{}).Where("id = ?", id).Update("imagen", categoria.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			return
		}
		categoria.Imagen = utils.SERVIMG + categoria.Imagen

	} else {
		categoria.Imagen = utils.SERVIMG + "default_user.png"
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Categoria actualizada correctamente", c, http.StatusOK)
}

func DeleteCategoriaDonacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.CategoriaDonacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar categoria"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Categoria eliminada exitosamente", c, http.StatusOK)
}
