package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPagoMarkets(c *gin.Context) {
	pagos := []*models.PagoMarket{}
	err := models.Db.Find(&pagos).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener pagos"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, pagos, c, http.StatusOK)
}

func GetPagoMarketPorId(c *gin.Context) {
	pago := &models.PagoMarket{}
	id := c.Param("id")
	err := models.Db.First(pago, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("PagoMarket no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener pago"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, pago, c, http.StatusOK)
}

func CreatePagoMarket(c *gin.Context) {
	pago := &models.PagoMarket{}
	err := c.ShouldBindJSON(pago)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Create(pago).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear pago"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "PagoMarket creada correctamente", c, http.StatusCreated)

}

func UpdatePagoMarket(c *gin.Context) {
	pago := &models.PagoMarket{}

	err := c.ShouldBindJSON(pago)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(pago).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar pago"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "PagoMarket actualizada correctamente", c, http.StatusOK)
}

func DeletePagoMarket(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.PagoMarket{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar pago"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "PagoMarket eliminada exitosamente", c, http.StatusOK)
}
