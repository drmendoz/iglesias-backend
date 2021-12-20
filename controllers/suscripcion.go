package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
)

func CrearSuscripcion(c *gin.Context) {
	idResidente := uint(c.GetInt("id_residente"))
	suscripcion := &models.Suscripcion{}
	var err error
	// err = c.ShouldBindJSON(suscripcion)
	// if err != nil {
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Parametros invalidos"), nil, c, http.StatusBadRequest)
	// 	return
	// }
	suscripcion.FechaInicio = time.Now().In(tiempo.Local)
	suscripcion.FechaFin = suscripcion.FechaInicio.AddDate(0, 1, 0)
	suscripcion.ResidenteID = idResidente
	suscripcion.Continua = true

	// Validar Pago de suscripcion

	err = models.Db.Create(suscripcion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear suscripcion"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, "Suscripcion creada exitosamente", c, http.StatusInternalServerError)
}

func AnularSuscripcion(c *gin.Context) {
	idResidente := uint(c.GetInt("id_residente"))
	err := models.Db.Where("residente_id = ?", idResidente).Updates(&models.Suscripcion{Continua: false}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al anular suscripcion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Suscripcion anulada exitosamente", c, http.StatusOK)

}

func RenovarSuscripciones(c *gin.Context) {
	fechaActual := time.Now().In(tiempo.Local)
	fechActualInicio := time.Date(fechaActual.Year(), fechaActual.Month(), fechaActual.Day(), 0, 0, 0, 0, time.Local)
	fechaActualFin := time.Date(fechaActual.Year(), fechaActual.Month(), fechaActual.Day(), 23, 59, 59, 59, time.Local)
	susPorVencer := []*models.Suscripcion{}
	err := models.Db.Where("continua = ?", true).Where("fecha_fin > ?", fechActualInicio).Where("fecha_fin < ?", fechaActualFin).Find(&susPorVencer).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al encontrar suscripciones activas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, sus := range susPorVencer {
		//Falta funcion de cobrar aqui
		fechaFin := sus.FechaFin.AddDate(0, 1, 0)
		err = models.Db.Create(&models.Suscripcion{FechaInicio: sus.FechaFin, FechaFin: fechaFin, Continua: true, ResidenteID: sus.ResidenteID}).Error
		if err != nil {
			_ = c.Error(err)
			errLog := fmt.Sprintf("Error al renovar suscripcion %d", sus.ID)
			_ = c.Error(errors.New(errLog))

		}
	}
	utils.CrearRespuesta(nil, "Suscripciones renovadas", c, http.StatusOK)

}

type SuscripcionResponse struct {
	Suscrito               bool  `json:"suscrito"`
	EmprendimientosActivos int64 `json:"emprendimientos_activos"`
}

func VerificarSuscripcionResidente(c *gin.Context) {
	idResidente := uint(c.GetInt("id_residente"))
	//suscrito, err := verificarSuscripcion(idResidente)
	// if err != nil {
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Error al verificar suscripcion"), nil, c, http.StatusInternalServerError)
	// 	return
	// }
	fechaActual := time.Now().In(tiempo.Local)
	var numEmp int64
	err := models.Db.Model(&models.Emprendimiento{}).Where("fecha_publicacion < ?", fechaActual).Where("fecha_vencimiento > ?", fechaActual).Where("residente_id = ?", idResidente).Count(&numEmp).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener publicacion"), nil, c, http.StatusInternalServerError)
		return
	}
	suscripcion := &SuscripcionResponse{Suscrito: true, EmprendimientosActivos: numEmp}
	utils.CrearRespuesta(nil, suscripcion, c, http.StatusOK)
}

func verificarSuscripcion(idResidente uint) (bool, error) {
	fechaActual := time.Now().In(tiempo.Local)
	var count int64
	err := models.Db.Model(&models.Suscripcion{}).Where("fecha_inicio < ?", fechaActual).Where("fecha_fin > ?", fechaActual).Where("residente_id = ?", idResidente).Count(&count).Error

	suscrito := true
	if count == 0 {
		suscrito = false
	}
	return suscrito, err
}

func ObtenerSuscripcionesResidente(c *gin.Context) {

}
