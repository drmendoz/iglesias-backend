package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func GetVentaById(c *gin.Context) {
	venta := &models.Venta{}
	err := models.Db.Find(venta, c.Param("id")).Error
	fmt.Println(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener modulo venta"), nil, c, http.StatusInternalServerError)
		return
	}
	if venta.Imagen != "" {
		venta.Imagen = utils.SERVIMG + venta.Imagen
	} else {
		venta.Imagen = utils.DefaultEtapa
	}
	utils.CrearRespuesta(err, venta, c, http.StatusOK)
}

func GetAllVentas(c *gin.Context) {
	venta := []*models.Venta{}
	err := models.Db.Find(&venta).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener ventas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, venta := range venta {
		if venta.Imagen == "" {
			venta.Imagen = utils.SERVIMG + venta.Imagen
		} else {
			venta.Imagen = utils.DefaultEtapa
		}
	}
	utils.CrearRespuesta(err, venta, c, http.StatusOK)
}

func CreateVenta(c *gin.Context) {
	venta := &models.Venta{}
	err := c.ShouldBindJSON(venta)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear venta"), nil, c, http.StatusBadRequest)
		return
	}
	txn := models.Db.Begin()
	err = txn.Omit("imagen").Create(venta).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear venta"), nil, c, http.StatusInternalServerError)
		return
	}
	txn.Commit()
	utils.CrearRespuesta(err, "Venta creada exitosamente", c, http.StatusCreated)
}

func UpdateVenta(c *gin.Context) {
	venta := &models.Venta{}
	err := c.ShouldBindJSON(venta)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(venta).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar Venta"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Venta actualizado exitosamente", c, http.StatusOK)
}

func DeleteVenta(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Venta{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar modulo venta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Modulo Venta eliminado exitosamente", c, http.StatusOK)
}
