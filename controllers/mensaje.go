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

func GetMensajes(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	mensajes := []*models.Mensaje{}
	err := models.Db.Order("created_at DESC").Where("etapa_id = ?", uint(idParroquia)).Preload("Respuestas", func(db *gorm.DB) *gorm.DB {
		return db.Order("respuesta_mensaje.created_at ASC")
	}).Preload("Noticia").Joins("Autor").Find(&mensajes).Error

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener mensajes"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, vot := range mensajes {
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
	utils.CrearRespuesta(nil, mensajes, c, http.StatusOK)
}

func GetMensajePorId(c *gin.Context) {
	mensaje := &models.Mensaje{}
	id := c.Param("id")
	err := models.Db.Preload("Respuestas", func(db *gorm.DB) *gorm.DB {
		return db.Order("respuesta_mensaje.created_at ASC")
	}).Preload("Autor", func(db *gorm.DB) *gorm.DB {
		return db.Joins("Fiel")
	}).Preload("Noticia").First(mensaje, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Mensaje no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener mensaje"), nil, c, http.StatusInternalServerError)
		return
	}

	imagenesArr := strings.Split(mensaje.Imagenes, ",")
	imagenes := []string{}
	if mensaje.Imagenes != "" {
		for _, imagen := range imagenesArr {
			imagen = utils.SERVIMG + imagen
			imagenes = append(imagenes, imagen)
		}
	} else {
		imagenes = append(imagenes, utils.DefaultMensaje)
	}
	mensaje.ImagenesArray = imagenes
	utils.CrearRespuesta(nil, mensaje, c, http.StatusOK)
}

func CreateMensaje(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	idUsuario := c.GetInt("id_usuario")
	if idParroquia == 0 {
		utils.CrearRespuesta(errors.New("No existe el id_etapa"), nil, c, http.StatusOK)
		return
	}

	mensaje := &models.Mensaje{}
	err := c.ShouldBindJSON(mensaje)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	mensaje.ParroquiaID = uint(idParroquia)
	mensaje.AutorID = uint(idUsuario)

	imagenesArr := mensaje.ImagenesArray
	if len(imagenesArr) > 0 {
		idUrb := fmt.Sprintf("%d", mensaje.ParroquiaID)
		imagenes := []string{}
		for _, imagen := range imagenesArr {
			imagen, err = img.FromBase64ToImage(imagen, "mensajes/"+time.Now().Format(time.RFC3339Nano)+idUrb, false)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear mensaje "), nil, c, http.StatusInternalServerError)
				return
			}
			imagenes = append(imagenes, imagen)
		}
		mensaje.Imagenes = strings.Join(imagenes, ",")
	}

	if mensaje.NoticiaID != nil {
		tx.Where("id = ?", mensaje.NoticiaID).Updates(&models.Publicacion{Adjuntado: true})
	}

	err = tx.Create(mensaje).Error
	tx.Commit()
	utils.CrearRespuesta(err, "Mensaje enviado exitosamente, prontro el administrador le responderá", c, http.StatusCreated)

}

func UpdateMensaje(c *gin.Context) {
	mensaje := &models.Mensaje{}

	err := c.ShouldBindJSON(mensaje)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("Opciones", "imagen").Where("id = ?", id).Updates(mensaje).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar mensaje"), nil, c, http.StatusInternalServerError)
		return
	}
	imagenesArr := mensaje.ImagenesArray
	if len(imagenesArr) > 0 {
		idUrb := fmt.Sprintf("%d", mensaje.ParroquiaID)
		imagenes := []string{}
		for _, imagen := range imagenesArr {
			imagen, err = img.FromBase64ToImage(imagen, "mensajes/"+time.Now().Format(time.RFC3339Nano)+idUrb, false)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear mensaje "), nil, c, http.StatusInternalServerError)
				return
			}
			imagenes = append(imagenes, imagen)
		}
		mensaje.Imagenes = strings.Join(imagenes, ",")
		err = tx.Model(&models.Mensaje{}).Where("id = ?", id).Updates(mensaje).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear mensaje "), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Mensaje actualizado correctamente", c, http.StatusOK)
}

func ActualizarMensajeRespuestas(c *gin.Context) {
	respuesta := &models.RespuestaMensaje{}
	err := c.ShouldBindJSON(respuesta)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagenes_string").Where("mensaje_id = ?", id).Updates(respuesta).Error
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Mensaje actualizado correctamente", c, http.StatusOK)
}

func DeleteMensaje(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Mensaje{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar mensaje"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Mensaje eliminado exitosamente", c, http.StatusOK)
}
