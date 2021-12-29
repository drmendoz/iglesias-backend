package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetActividads(c *gin.Context) {
	idParroquia := c.GetInt("id_parroquia")
	etps := []*models.Actividad{}
	err := models.Db.Where(&models.Actividad{ParroquiaID: uint(idParroquia)}).Order("created_at ASC").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener misas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, act := range etps {
		act.Imagen = utils.SERVIMG + act.Imagen
		act.Video = utils.SERVIMG + act.Video
	}
	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetActividadPorId(c *gin.Context) {
	etp := &models.Actividad{}
	id := c.Param("id")
	err := models.Db.First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Actividad no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener misa"), nil, c, http.StatusInternalServerError)
		return
	}
	etp.Imagen = utils.SERVIMG + etp.Imagen
	etp.Video = utils.SERVIMG + etp.Video
	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateActividad(c *gin.Context) {
	etp := &models.Actividad{}
	idParroquia := c.GetInt("id_parroquia")

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = uint(idParroquia)

	tx := models.Db.Begin()
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear misa"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Actividad creada correctamente", c, http.StatusCreated)

}

func UpdateActividad(c *gin.Context) {
	etp := &models.Actividad{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	if strings.HasPrefix(etp.Imagen, "https://") {
		etp.Imagen = ""
	}
	if strings.HasPrefix(etp.Video, "https://") {
		etp.Video = ""
	}
	tx := models.Db.Begin()
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar misa"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Actividad actualizada correctamente", c, http.StatusOK)
}

func DeleteActividad(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Actividad{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar misa"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Actividad eliminada exitosamente", c, http.StatusOK)
}
