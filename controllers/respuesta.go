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

func GetRespuestas(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	respuestas := []*models.RespuestaMensaje{}
	var err error
	err = models.Db.Order("created_at DESC").Where("etapa_id = ?", idParroquia).Joins("Mensaje").Joins("Autor").Find(&respuestas).Error

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener respuestas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, vot := range respuestas {
		imagenesArr := strings.Split(vot.Imagenes, ",")
		imagenes := []string{}
		if vot.Imagenes != "" {
			for _, imagen := range imagenesArr {
				imagen = utils.SERVIMG + imagen
				imagenes = append(imagenes, imagen)
			}
		} else {
			imagenes = append(imagenes, utils.DefaultMensaje)
		}
		vot.ImagenesArray = imagenes
	}
	utils.CrearRespuesta(nil, respuestas, c, http.StatusOK)
}

func GetRespuestaPorId(c *gin.Context) {
	respuesta := &models.RespuestaMensaje{}
	id := c.Param("id")
	err := models.Db.Joins("Mensaje").Joins("Autor").First(respuesta, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Respuesta no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener respuesta"), nil, c, http.StatusInternalServerError)
		return
	}

	imagenesArr := strings.Split(respuesta.Imagenes, ",")
	imagenes := []string{}
	if respuesta.Imagenes != "" {
		for _, imagen := range imagenesArr {
			imagen = utils.SERVIMG + imagen
			imagenes = append(imagenes, imagen)
		}
	} else {
		imagenes = append(imagenes, utils.DefaultMensaje)
	}
	respuesta.ImagenesArray = imagenes
	utils.CrearRespuesta(nil, respuesta, c, http.StatusOK)
}

func CreateRespuesta(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	idUsuario := c.GetInt("id_usuario")
	if idParroquia == 0 {
		utils.CrearRespuesta(errors.New("No existe el id_etapa"), nil, c, http.StatusOK)
		return
	}

	respuesta := &models.RespuestaMensaje{}
	err := c.ShouldBindJSON(respuesta)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	respuesta.AutorID = uint(idUsuario)

	println("respuesta.AutorID")
	println(respuesta.AutorID)

	imagenesArr := respuesta.ImagenesArray
	if len(imagenesArr) > 0 {
		idUrb := fmt.Sprintf("%d", idParroquia)
		imagenes := []string{}
		for _, imagen := range imagenesArr {
			imagen, err = img.FromBase64ToImage(imagen, "respuestas/"+time.Now().Format(time.RFC3339Nano)+idUrb, false)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear respuesta "), nil, c, http.StatusInternalServerError)
				return
			}
			imagenes = append(imagenes, imagen)
		}
		respuesta.Imagenes = strings.Join(imagenes, ",")
	}

	err = tx.Create(respuesta).Error
	tx.Commit()
	utils.CrearRespuesta(err, "Respuesta creado correctamente", c, http.StatusCreated)

}

func UpdateRespuesta(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	respuesta := &models.RespuestaMensaje{}

	err := c.ShouldBindJSON(respuesta)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagenes_string").Where("id = ?", id).Updates(respuesta).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar respuesta"), nil, c, http.StatusInternalServerError)
		return
	}
	imagenesArr := respuesta.ImagenesArray
	if len(imagenesArr) > 0 {
		idUrb := fmt.Sprintf("%d", idParroquia)
		imagenes := []string{}
		for _, imagen := range imagenesArr {
			imagen, err = img.FromBase64ToImage(imagen, "respuestas/"+time.Now().Format(time.RFC3339Nano)+idUrb, false)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear respuesta "), nil, c, http.StatusInternalServerError)
				return
			}
			imagenes = append(imagenes, imagen)
		}
		err = tx.Model(&models.RespuestaMensaje{}).Where("id = ?", id).Update("imagenes_string", strings.Join(imagenes, ",")).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear respuesta "), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Respuesta actualizada correctamente", c, http.StatusOK)
}

func DeleteRespuesta(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.RespuestaMensaje{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar respuesta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Respuesta eliminada exitosamente", c, http.StatusOK)
}
