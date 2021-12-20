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

func GetContactoEtapas(c *gin.Context) {
	contactos := []*models.ContactoEtapa{}
	var err error
	idEtapa := c.GetInt("id_etapa")
	if idEtapa != 0 {
		err = models.Db.Where("etapa_id = ?", idEtapa).Find(&contactos).Error
	} else {
		err = models.Db.Find(&contactos).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener contactos"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, contacto := range contactos {
		if contacto.Imagen == "" {
			contacto.Imagen = utils.DefaultNoticia
		} else {
			contacto.Imagen = utils.SERVIMG + contacto.Imagen
		}
	}
	utils.CrearRespuesta(err, contactos, c, http.StatusOK)
}

func GetContactoEtapaPorId(c *gin.Context) {
	contacto := &models.ContactoEtapa{}
	id := c.Param("id")
	err := models.Db.First(contacto, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Contacto etapa no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener contacto"), nil, c, http.StatusInternalServerError)
		return
	}
	if contacto.Imagen == "" {
		contacto.Imagen = utils.DefaultNoticia
	} else {
		contacto.Imagen = utils.SERVIMG + contacto.Imagen
	}

	utils.CrearRespuesta(nil, contacto, c, http.StatusOK)
}

func CreateContactoEtapa(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))

	contacto := &models.ContactoEtapa{}
	err := c.ShouldBindJSON(contacto)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	contacto.EtapaID = idEtapa
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(contacto).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear contacto"), nil, c, http.StatusInternalServerError)
		return
	}

	if contacto.Imagen == "" {
		contacto.Imagen = utils.DefaultNoticia
	} else {
		idUrb := fmt.Sprintf("%d", contacto.ID)
		contacto.Imagen, err = img.FromBase64ToImage(contacto.Imagen, "contactos/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(contacto.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear contacto "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.ContactoEtapa{}).Where("id = ?", contacto.ID).Update("imagen", contacto.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear contacto "), nil, c, http.StatusInternalServerError)
			return
		}
		contacto.Imagen = utils.SERVIMG + contacto.Imagen
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Contacto etapa creada exitosamente", c, http.StatusCreated)

}

func UpdateContactoEtapa(c *gin.Context) {
	contacto := &models.ContactoEtapa{}

	err := c.ShouldBindJSON(contacto)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(contacto).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar contacto"), nil, c, http.StatusInternalServerError)
		return
	}
	if img.IsBase64(contacto.Imagen) {
		idUrb := fmt.Sprintf("%d", contacto.ID)
		contacto.Imagen, err = img.FromBase64ToImage(contacto.Imagen, "contactos/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear contacto "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.ContactoEtapa{}).Where("id = ?", id).Update("imagen", contacto.Imagen).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar contacto"), nil, c, http.StatusInternalServerError)
			return
		}
		contacto.Imagen = utils.SERVIMG + contacto.Imagen

	} else {
		contacto.Imagen = utils.DefaultNoticia
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Contacto etapa actualizada correctamente", c, http.StatusOK)
}

func DeleteContactoEtapa(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.ContactoEtapa{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar contacto"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Contacto etapa eliminada exitosamente", c, http.StatusOK)
}
