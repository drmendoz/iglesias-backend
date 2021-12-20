package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/paymentez"
	"github.com/gin-gonic/gin"
)

func GetTarjetas(c *gin.Context) {
	idFiel := c.GetInt("id_residente")
	tarjetas, err := paymentez.GetTarjetas(int64(idFiel))
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener tarjetas"), nil, c, 500)
		return
	}
	utils.CrearRespuesta(nil, tarjetas, c, http.StatusOK)
}

func DeleteTarjeta(c *gin.Context) {
	idFiel := c.GetInt("id_residente")
	tokenTarjeta := c.Param("token")
	res, err := paymentez.DeleteTarjeta(idFiel, tokenTarjeta)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar tarjeta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, res, c, http.StatusOK)
}

func CobrarTarjeta(c *gin.Context) {
	idFiel := c.GetInt("id_residente")
	tokenTarjeta := c.Param("token")
	id := fmt.Sprintf("%d", idFiel)
	res, err := paymentez.CobrarTarjeta(id, "drmendozal98@gmail.com", 112, "Prueba", "2", 12, tokenTarjeta)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar tarjeta"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, res, c, http.StatusOK)
}
