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
)

func GetParroquias(c *gin.Context) {
	etps := []*models.Parroquia{}
	err := models.Db.Order("Nombre ASC").Preload("Iglesia").Preload("Modulos").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener parroquias"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, etp := range etps {
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultParroquia
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}

	}
	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetParroquiaPorId(c *gin.Context) {
	etp := &models.Parroquia{}
	id := c.Param("id")
	err := models.Db.Preload("Iglesia").Preload("Modulos").First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Parroquia no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultParroquia
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateParroquia(c *gin.Context) {
	etp := &models.Parroquia{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear etapa"), nil, c, http.StatusInternalServerError)
		return
	}

	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultParroquia
	} else {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "parroquias/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(etp.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.ModulosParroquia{}).Create(&models.ModulosParroquia{ParroquiaID: etp.ID, Horario: true, Actividad: true, Emprendimiento: true, Intencion: true, Musica: true, Ayudemos: true, Curso: true, Matrimonio: true, Galeria: true, Publicacion: true}).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)
			return
		}
		err = tx.Model(&models.Parroquia{}).Where("id = ?", etp.ID).Update("imagen", etp.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)
			return
		}
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Parroquia creada correctamente", c, http.StatusCreated)

}

func UpdateParroquia(c *gin.Context) {
	etp := &models.Parroquia{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	idPar, _ := strconv.Atoi(id)
	etp.ID = uint(idPar)
	tx := models.Db.Begin()

	err = tx.Omit("imagen").Where("id = ?", id).Updates(etp).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	err = tx.Model(&models.Parroquia{}).Where("id = ?", id).Updates(map[string]interface{}{
		"boton_pago_matrimonio":     etp.BotonPagoMatrimonio,
		"boton_pago_emprendimiento": etp.BotonPagoEmprendimiento,
		"boton_pago_curso":          etp.BotonPagoCurso,
		"boton_pago_intencion":      etp.BotonPagoIntencion,
		"boton_pago_actividad":      etp.BotonPagoActividad,
		"boton_pago_musica":         etp.BotonPagoMusica,
	}).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
		return
	}

	if etp.Imagen != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "parroquias/"+time.RFC3339+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Parroquia{}).Where("id = ?", etp.ID).Update("imagen", etp.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
			return
		}
		etp.Imagen = utils.SERVIMG + etp.Imagen

	} else {
		etp.Imagen = utils.DefaultParroquia
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Parroquia actualizada correctamente", c, http.StatusOK)
}

func UpdateModulosParroquia(c *gin.Context) {
	etp := &models.ModulosParroquia{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")

	tx := models.Db.Begin()

	err = tx.Model(&models.ModulosParroquia{}).Where("parroquia_id = ?", id).Updates(map[string]interface{}{
		"horario":        etp.Horario,
		"actividad":      etp.Actividad,
		"emprendimiento": etp.Emprendimiento,
		"intencion":      etp.Intencion,
		"musica":         etp.Musica,
		"ayudemos":       etp.Ayudemos,
		"galeria":        etp.Galeria,
		"matrimonio":     etp.Matrimonio,
		"publicacion":    etp.Publicacion,
		"curso":          etp.Curso,
	}).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Parroquia actualizada correctamente", c, http.StatusOK)
}

func DeleteParroquia(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Parroquia{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Parroquia eliminada exitosamente", c, http.StatusOK)
}
