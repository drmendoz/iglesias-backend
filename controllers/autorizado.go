package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAutorizados(c *gin.Context) {
	autorizados := []*models.Autorizado{}

	err := models.Db.Limit(100).Where(&models.Autorizado{TipoUsuario: "EXPRESO"}).Order("created_at desc").Preload("Autorizaciones", func(db *gorm.DB) *gorm.DB {
		return db.Joins("Casa").Order("Casa.Manzana ASC").Order("Casa.Manzana ASC")
	}).Find(&autorizados).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener autorizados"), nil, c, http.StatusInternalServerError)
		return
	}

	if len(autorizados) > 0 {
		for _, autorizado := range autorizados {
			if autorizado.Imagen != "" {
				if !strings.HasPrefix(autorizado.Imagen, "https://") {
					autorizado.Imagen = utils.SERVIMG + autorizado.Imagen
				}
			} else {
				autorizado.Imagen = utils.DefaultExpreso
			}
			if autorizado.Pdf != "" {
				autorizado.Pdf = "https://api.practical.com.ec/public/pdf/" + autorizado.Pdf
			}
		}
	}

	utils.CrearRespuesta(err, autorizados, c, http.StatusOK)
}

func GetAutorizadoPorId(c *gin.Context) {
	autorizado := &models.Autorizado{}
	id := c.Param("id")
	err := models.Db.First(autorizado, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Autorizado no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener autorizado"), nil, c, http.StatusInternalServerError)
		return
	}
	if autorizado.Imagen == "" {
		autorizado.Imagen = utils.DefaultVisita
	} else {
		autorizado.Imagen = utils.SERVIMG + autorizado.Imagen
	}
	utils.CrearRespuesta(nil, autorizado, c, http.StatusOK)
}

func CreateAutorizado(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	idParroquia := c.GetInt("id_etapa")
	if idUsuario == 0 || idParroquia == 0 {
		utils.CrearRespuesta(errors.New("Error al crear autorizado"), nil, c, http.StatusInternalServerError)
		return
	}
	autorizado := &models.Autorizado{}
	err := c.ShouldBindJSON(autorizado)
	if autorizado.Pdf != "" {
		uri := strings.Split(autorizado.Pdf, ";")[0]
		if uri == "data:application/pdf" {
			nombre := fmt.Sprintf("autorizacion-%d.pdf", time.Now().Unix())
			base64 := strings.Split(autorizado.Pdf, ",")[1]
			err = utils.SubirPdf(nombre, base64)
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
				return
			}
			autorizado.Pdf = nombre
		} else {
			autorizado.Pdf = ""
		}
	} else {
		autorizado.Pdf = ""
	}
	autorizado.PublicadorID = uint(idUsuario)
	autorizado.ParroquiaID = uint(idParroquia)

	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear autorizado"), nil, c, http.StatusInternalServerError)
		return
	}

	if autorizado.Imagen != "" {
		autorizado.Imagen, err = img.FromBase64ToImage(autorizado.Imagen, "autorizados/"+time.Now().Format(time.RFC3339), false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear autorizado "), nil, c, http.StatusInternalServerError)

			return
		}
	}

	err = tx.Create(autorizado).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(err, "Error al crear autorización", c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Autorizado creado correctamente", c, http.StatusCreated)
}

func UpdateAutorizado(c *gin.Context) {
	autorizado := &models.Autorizado{}

	err := c.ShouldBindJSON(autorizado)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	autorizado.UpdatedAt = time.Now()
	if autorizado.Pdf != "" {
		uri := strings.Split(autorizado.Pdf, ";")[0]
		if uri == "data:application/pdf" {
			nombre := fmt.Sprintf("autorizacion-%d.pdf", time.Now().Unix())
			base64 := strings.Split(autorizado.Pdf, ",")[1]
			err = utils.SubirPdf(nombre, base64)
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
				return
			}
			autorizado.Pdf = nombre
		} else {
			autorizado.Pdf = ""
		}
	} else {
		autorizado.Pdf = ""
	}
	err = tx.Omit("imagen").Where("id = ?", id).Updates(autorizado).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar autorizado"), nil, c, http.StatusInternalServerError)
		return
	}
	if autorizado.Imagen != "" {
		idUrb := fmt.Sprintf("%d", autorizado.ID)
		autorizado.Imagen, err = img.FromBase64ToImage(autorizado.Imagen, "autorizados/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear autorizado "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Autorizado{}).Where("id = ?", id).Update("imagen", autorizado.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar autorizado"), nil, c, http.StatusInternalServerError)
			return
		}
		autorizado.Imagen = utils.SERVIMG + autorizado.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Autorizado actualizada correctamente", c, http.StatusOK)
}

func DeleteAutorizado(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Autorizado{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar autorizado"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Autorizado eliminada exitosamente", c, http.StatusOK)
}
