package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

type Notificaciones struct {
	Emprendimiento  int `json:"emprendimientos"`
	Bitacora        int `json:"bitacora"`
	Buzon           int `json:"buzon"`
	Galeria         int `json:"galeria"`
	Camara          int `json:"camaras"`
	Alicuota        int `json:"alicuota"`
	AreaSocial      int `json:"area-social"`
	Administradores int `json:"administrador"`
	Votacion        int `json:"votacion"`
	Reserva         int `json:"reserva"`
}

func obtenerNotificaciones(idResidente int, idCasa int, idEtapa int) (*Notificaciones, error) {

	notificacion := &Notificaciones{}
	residente := &models.Residente{}
	err := models.Db.Select("visualizacion_emprendimiento", "visualizacion_galeria", "visualizacion_bitacora", "visualizacion_buzon", "visualizacion_camara", "visualizacion_administradores", "visualizacion_alicuota", "visualizacion_area_social", "visualizacion_votacion", "visualizacion_reservas").First(residente, idResidente).Error

	if err != nil {
		return nil, err
	}
	var countEmprendimiento int64
	var countBitacora int64
	var countBuzon int64
	var countGaleria int64
	var countCamara int64
	var countAlicouta int64
	var countAreaSocial int64
	var countAdministradores int64
	var countVotacion int64
	var countReservas int64
	err = models.Db.Model(&models.Visita{}).Where("created_at > ? and casa_id = ?", residente.VisualizacionBitacora, idCasa).Count(&countBitacora).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.Emprendimiento{}).Where("created_at > ? and fecha_vencimiento < ?", residente.VisualizacionEmprendimiento, residente.VisualizacionEmprendimiento).Count(&countEmprendimiento).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.Publicacion{}).Where("created_at > ? and etapa_id = ?", residente.VisualizacionBuzon, idEtapa).Count(&countBuzon).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.ImagenGaleria{}).Where("created_at > ? and etapa_id = ?", residente.VisualizacionGaleria, idEtapa).Count(&countGaleria).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.Votacion{}).Where("created_at > ? and etapa_id = ?", residente.VisualizacionVotacion, idEtapa).Count(&countVotacion).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.EtapaCamara{}).Where("created_at > ? and etapa_id = ?", residente.VisualizacionCamara, idEtapa).Count(&countCamara).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.Alicuota{}).Where("created_at > ? and casa_id = ?", residente.VisualizacionAlicuota, idCasa).Count(&countAlicouta).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.AreaSocial{}).Where("created_at > ? and etapa_id = ?", residente.VisualizacionAreaSocial, idEtapa).Count(&countAreaSocial).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.Administrativo{}).Where("created_at > ? and etapa_id = ?", residente.VisualizacionAdministradores, idEtapa).Count(&countAdministradores).Error
	if err != nil {
		return nil, err
	}
	err = models.Db.Model(&models.ReservacionAreaSocial{}).Where("created_at > ? and residente_id = ?", residente.VisualizacionReservas, idResidente).Count(&countReservas).Error
	if err != nil {
		return nil, err
	}
	notificacion.Bitacora = int(countBitacora)
	notificacion.Buzon = int(countBuzon)
	notificacion.Emprendimiento = int(countEmprendimiento)
	notificacion.Galeria = int(countGaleria)
	notificacion.Votacion = int(countVotacion)
	notificacion.Camara = int(countCamara)
	notificacion.Alicuota = int(countAlicouta)
	notificacion.AreaSocial = int(countAreaSocial)
	notificacion.Administradores = int(countAdministradores)
	notificacion.Reserva = int(countReservas)
	return notificacion, nil
}

func GetNotificacionesRequest(c *gin.Context) {
	idEtapa := c.GetInt("id_etapa")
	idResidente := c.GetInt("id_residente")
	idCasa := c.GetInt("id_casa")
	notificaciones, err := obtenerNotificaciones(idResidente, idCasa, idEtapa)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al obtener notificaciones"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, notificaciones, c, http.StatusOK)

}
