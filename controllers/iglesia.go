package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetIglesiaes(c *gin.Context) {
	igls := []*models.Iglesia{}
	err := models.Db.Order("nombre ASC").Find(&igls).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener iglanizaciones"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, igls, c, http.StatusOK)
}

func GetIglesiaPorId(c *gin.Context) {

	igl := &models.Iglesia{}
	id := c.GetInt("id_iglesia")
	if id == 0 {

		id, _ = strconv.Atoi(c.Param("id"))
	}

	err := models.Db.Preload("Etapas").First(igl, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Iglesia no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener iglanizacion"), nil, c, http.StatusInternalServerError)
		return
	}

	for _, etp := range igl.Parroquias {
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultParroquia
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}
	}
	utils.CrearRespuesta(nil, igl, c, http.StatusOK)
}

func CreateIglesia(c *gin.Context) {
	igl := &models.Iglesia{}
	err := c.ShouldBindJSON(igl)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Create(igl).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear iglanizacion"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Iglesia creada correctamente", c, http.StatusCreated)

}

func UpdateIglesia(c *gin.Context) {
	igl := &models.Iglesia{}

	err := c.ShouldBindJSON(igl)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(igl).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar iglanizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, "Iglesia actualizada correctamente", c, http.StatusOK)
}

func DeleteIglesia(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Iglesia{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar iglanizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Iglesia eliminada exitosamente", c, http.StatusOK)
}

func GetParroquiasPorIglesia(c *gin.Context) {
	igl := &models.Iglesia{}
	id := c.Param("id")
	err := models.Db.Preload(clause.Associations).First(igl, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Iglesia no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener iglanizacion"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, igl, c, http.StatusOK)
}

func GetIglesiaesCount(c *gin.Context) {
	var igls int64
	err := models.Db.Model(&models.Iglesia{}).Count(&igls).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, igls, c, http.StatusOK)
}
