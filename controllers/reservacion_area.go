package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/paymentez"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetReservacionAreaSocials(c *gin.Context) {
	reservas := []*models.ReservacionAreaSocial{}
	err := models.Db.Order("hora_inicio DESC").Find(&reservas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener reservas"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, reservas, c, http.StatusOK)
}

func GetReservacionAreaSocialPorId(c *gin.Context) {
	reserva := &models.ReservacionAreaSocial{}
	id := c.Param("id")
	err := models.Db.First(reserva, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("ReservacionAreaSocial no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener reserva"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, reserva, c, http.StatusOK)
}

func CreateReservacionAreaSocial(c *gin.Context) {
	reserva := &models.ReservacionAreaSocial{}
	err := c.ShouldBindJSON(reserva)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Create(reserva).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear reserva"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "ReservacionAreaSocial creada correctamente", c, http.StatusCreated)

}

func UpdateReservacionAreaSocial(c *gin.Context) {
	reserva := &models.ReservacionAreaSocial{}

	err := c.ShouldBindJSON(reserva)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(reserva).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar reserva"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "ReservacionAreaSocial actualizada correctamente", c, http.StatusOK)
}

func DeleteReservacionAreaSocial(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.ReservacionAreaSocial{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar reserva"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "ReservacionAreaSocial eliminada exitosamente", c, http.StatusOK)
}

func GetReservacionesFielAreaSocial(c *gin.Context) {
	idFiel := c.GetInt("id_residente")
	reservas := []*models.ReservacionAreaSocial{}
	err := models.Db.Order("created_at DESC").Where("residente_id = ?", idFiel).Joins("AreaSocial").Find(&reservas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener reservas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, res := range reservas {
		if res.AreaSocial.Imagen == "" {
			res.AreaSocial.Imagen = utils.DefaultAreaSocial
		} else {
			res.AreaSocial.Imagen = utils.SERVIMG + res.AreaSocial.Imagen
		}
		if res.AreaSocial.ImagenReserva == "" {
			res.AreaSocial.ImagenReserva = utils.DefaultAreaSocial
		} else {
			res.AreaSocial.ImagenReserva = utils.SERVIMG + res.AreaSocial.ImagenReserva
		}
	}
	utils.CrearRespuesta(err, reservas, c, http.StatusOK)
}

type PagoReserva struct {
	HoraInicio   *time.Time `json:"hora_inicio"`
	HoraFin      *time.Time `json:"hora_fin"`
	TokenTarjeta string     `json:"token_tarjeta"`
}

func CreateReservaFiel(c *gin.Context) {
	pago := &PagoReserva{}
	idFiel := c.GetInt("id_residente")
	idArea := c.Param("id")
	err := c.ShouldBindJSON(pago)
	areaId, err := strconv.Atoi(idArea)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error de parametros"), nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	res := &models.Fiel{}
	err = tx.Joins("Usuario").Find(res, idFiel).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	idRes := fmt.Sprintf("%d", idFiel)
	reserva := &models.ReservacionAreaSocial{}
	reserva.HoraInicio = *pago.HoraInicio
	reserva.HoraFin = *pago.HoraFin
	reserva.FielID = uint(idFiel)
	reserva.AreaSocialID = uint(areaId)
	area := &models.AreaSocial{}
	err = tx.First(area, areaId).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe el area social"), nil, c, http.StatusBadRequest)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener area social"), nil, c, http.StatusInternalServerError)
		return
	}

	var result int64
	inicioMes := time.Date(reserva.HoraInicio.Year(), reserva.HoraInicio.Month(), 1, 0, 0, 0, 0, tiempo.Local)
	finMes := time.Date(reserva.HoraInicio.Year(), reserva.HoraInicio.Month(), 30, 0, 0, 0, 0, tiempo.Local)
	err = tx.Model(&models.ReservacionAreaSocial{}).Where("hora_inicio between ? and ?", inicioMes, finMes).Where("residente_id = ?", idFiel).Where("area_social_id = ?", areaId).Count(&result).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear reservacion"), nil, c, http.StatusInternalServerError)
		return
	}
	if result >= int64(area.ReservasFielMes) {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Ha excedido el numero de reservas de esta area social en este mes"), nil, c, http.StatusForbidden)
		return
	}
	err = tx.Create(reserva).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear reserva"), nil, c, http.StatusInternalServerError)
		return
	}
	if area.Precio > 0 {
		tarjeta := &models.FielTarjeta{}
		err = tx.Where("token_tarjeta = ?", pago.TokenTarjeta).First(tarjeta).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tarjeta.TokenTarjeta = pago.TokenTarjeta
				tarjeta.FielID = uint(idFiel)
				err = tx.Create(&tarjeta).Error
				if err != nil {
					tx.Rollback()
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error con tarjeta"), nil, c, http.StatusInternalServerError)
					return
				}
			} else {

				tx.Rollback()
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
				return

			}
		}
		trans := &models.Transaccion{FielTarjetaID: tarjeta.ID, Tipo: "RES"}
		err = tx.Create(trans).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
			return
		}
		idTrans := fmt.Sprintf("%d", trans.ID)
		descripcion := fmt.Sprintf("Pago  de reserva area social # %d", reserva.ID)
		cobro, err := paymentez.CobrarTarjeta(idRes, res.Usuario.Correo, area.Precio, descripcion, idTrans, 0, pago.TokenTarjeta)
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al debitar tarjeta"), nil, c, http.StatusPaymentRequired)
			return
		}
		montoReal := fmt.Sprintf("%f", cobro.Transaccion.Monto)
		transNueva := &models.Transaccion{Estado: cobro.Transaccion.Status, DiaPago: cobro.Transaccion.FechaPago, Monto: montoReal, CodigoAutorizacion: cobro.Transaccion.CodigoAutorizacion, Mensaje: cobro.Transaccion.Mensaje, Descripcion: descripcion, FielTarjetaID: tarjeta.ID}
		err = tx.Where("id = ?", trans.ID).Updates(transNueva).Error
		if err != nil {
			_ = c.Error(err)
			_ = c.Error(errors.New("Reserva pagada pero error al guardar transaccion"))
			tx.Commit()
			utils.CrearRespuesta(nil, "Reserva pagada exitosamente", c, http.StatusOK)
			return
		}

		err = tx.Where("id= ?", reserva.ID).Updates(&models.ReservacionAreaSocial{TransaccionID: &trans.ID, ValorCancelado: area.Precio}).Error
		if err != nil {
			_ = c.Error(err)
		}
	}
	tx.Commit()
	idCasa := c.GetInt("id_casa")
	idParroquia := c.GetInt("id_etapa")
	visualizaciones, _ := obtenerNotificaciones(idFiel, idCasa, idParroquia)

	utils.CrearRespuesta(nil, &VisuMensajes{Notificaciones: visualizaciones, Mensaje: "Reservación creada exitosamente"}, c, http.StatusCreated)

}

type VisuMensajes struct {
	Notificaciones *Notificaciones `json:"notificaciones"`
	Mensaje        string          `json:"mensaje"`
}
