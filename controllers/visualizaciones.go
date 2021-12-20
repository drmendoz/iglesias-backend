package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func ActualizarVisualizacion(c *gin.Context) {
	idRes := c.GetInt("id_residente")
	fechaActual := time.Now()
	modulo := c.Param("modulo")
	var err error
	switch modulo {
	case "emprendimiento":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionEmprendimiento: &fechaActual}).Error
	case "bitacora":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionBitacora: &fechaActual}).Error
	case "galeria":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionGaleria: &fechaActual}).Error
	case "buzon":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionBuzon: &fechaActual}).Error
	case "votacion":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionVotacion: &fechaActual}).Error
	case "administrador":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionAdministradores: &fechaActual}).Error
	case "camara":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionCamara: &fechaActual}).Error
	case "area-social":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionAreaSocial: &fechaActual}).Error
	case "alicuota":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionAlicuota: &fechaActual}).Error
	case "reserva":
		err = models.Db.Where("id = ?", idRes).Updates(&models.Residente{VisualizacionReservas: &fechaActual}).Error
	}

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar visualizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	idEtapa := c.GetInt("id_etapa")
	idCasa := c.GetInt("id_casa")
	notificaciones, err := obtenerNotificaciones(idRes, idCasa, idEtapa)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al obtener notificaciones"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, notificaciones, c, http.StatusOK)
}
