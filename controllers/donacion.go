package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/paymentez"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDonacions(c *gin.Context) {
	etps := []*models.Donacion{}
	idParroquia := c.GetInt("id_parroquia")
	idCategoria := c.Query("id_categoria")
	idCat, err := strconv.Atoi(idCategoria)
	if err != nil {
		idCat = 0
	}
	err = models.Db.Where(&models.Donacion{ParroquiaID: uint(idParroquia), CategoriaDonacionID: uint(idCat)}).Order("created_at asc").Preload("CategoriaDonacion").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener areas sociales"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, etp := range etps {
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultDonacion
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}
		if etp.ImagenReserva == "" {
			etp.ImagenReserva = utils.DefaultDonacion
		} else {
			etp.ImagenReserva = utils.SERVIMG + etp.ImagenReserva
		}

	}
	utils.CrearRespuesta(nil, etps, c, http.StatusOK)
}

func GetDonacionPorID(c *gin.Context) {
	etp := &models.Donacion{}
	id := c.Param("id")
	err := models.Db.Preload("Aportaciones").Preload("Aportaciones.Fiel").Preload("Aportaciones.Transaccion").First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Doncacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultDonacion
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}
	if etp.ImagenReserva == "" {
		etp.ImagenReserva = utils.DefaultDonacion
	} else {
		etp.ImagenReserva = utils.SERVIMG + etp.ImagenReserva
	}
	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateDonacion(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_parroquia"))
	etp := &models.Donacion{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = idParroquia

	tx := models.Db.Begin()
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear donacion"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Donacion creada correctamente", c, http.StatusCreated)

}

func UpdateDonacion(c *gin.Context) {
	etp := &models.Donacion{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	if strings.HasPrefix(etp.Imagen, "https://") {
		etp.Imagen = ""
	}
	if strings.HasPrefix(etp.ImagenReserva, "https://") {
		etp.ImagenReserva = ""
	}
	tx := models.Db.Begin()
	id := c.Param("id")
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Donacion actualizada correctamente", c, http.StatusOK)
}

func DeleteDonacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Donacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Donacion eliminada exitosamente", c, http.StatusOK)
}

func AportarDonacion(c *gin.Context) {
	idFiel := c.GetInt("id_fiel")
	idParroquia := c.GetInt("id_parroquia")
	id := c.Param("id")
	idDon, err := strconv.Atoi(id)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al aportar donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	aportacion := &models.Aportacion{}
	err = c.ShouldBindJSON(aportacion)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al aportar donacion"), nil, c, http.StatusInternalServerError)
		return
	}
	if aportacion.TokenTarjeta == "" {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Ingrese tarjeta de credito"), nil, c, http.StatusNotAcceptable)
		return
	}
	aportacion.DonacionID = uint(idDon)
	aportacion.FielID = uint(idFiel)
	tx := models.Db.Begin()
	fiel := &models.Fiel{}
	err = tx.Joins("Usuario").Find(fiel, idFiel).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	donacion := &models.Donacion{}
	err = tx.First(donacion, idDon).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	tarjeta := &models.FielTarjeta{}
	err = tx.Where("token_tarjeta = ?", aportacion.TokenTarjeta).First(tarjeta).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tarjeta.TokenTarjeta = aportacion.TokenTarjeta
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
	aportacion.Transaccion = &models.Transaccion{FielTarjetaID: tarjeta.ID, CategoriaID: donacion.CategoriaDonacionID, ParroquiaID: uint(idParroquia), CasoID: uint(idDon)}
	err = tx.Create(aportacion).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	idTrans := fmt.Sprintf("%d", aportacion.Transaccion.ID)
	descripcion := fmt.Sprintf("Pago  de aportacion # %d", aportacion.ID)
	idFi := fmt.Sprintf("%d", idFiel)
	cobro, err := paymentez.CobrarTarjeta(idFi, fiel.Usuario.Correo, aportacion.Monto, descripcion, idTrans, 0, aportacion.TokenTarjeta)
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al debitar tarjeta"), nil, c, http.StatusPaymentRequired)
		return
	}
	montoReal := fmt.Sprintf("%f", cobro.Transaccion.Monto)
	transNueva := &models.Transaccion{Estado: cobro.Transaccion.Status, DiaPago: cobro.Transaccion.FechaPago, Monto: montoReal, CodigoAutorizacion: cobro.Transaccion.CodigoAutorizacion, Mensaje: cobro.Transaccion.Mensaje, Descripcion: descripcion, FielTarjetaID: tarjeta.ID}
	err = tx.Where("id = ?", aportacion.Transaccion.ID).Updates(transNueva).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(nil, "Aportacion creada exitosamente", c, http.StatusOK)
		return
	}
	_ = tx.Commit()
	utils.CrearRespuesta(nil, "Aportacion creada exitosamente", c, http.StatusOK)
}

type DonacionesTotal struct {
	Donaciones []*models.Aportacion `json:"aportaciones"`
	Monto      float64              `json:"monto"`
	Donacion   *models.Donacion     `json:"donacion"`
}

func GetAportacionesDeDonacion(c *gin.Context) {

	idDonacion := c.Param("id")
	idD, err := strconv.Atoi(idDonacion)
	if err != nil {
		utils.CrearRespuesta(errors.New("Id incorrecto"), nil, c, http.StatusBadRequest)
		return
	}
	donacion := &models.Donacion{}
	err = models.Db.Preload("CategoriaDonacion").First(donacion, idD).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener aportaciones"), nil, c, http.StatusInternalServerError)
		return
	}
	fechaInicioQ := c.Query("fecha_inicio")
	fechaFinQ := c.Query("fecha_fin")
	fechaInicio, err := time.Parse("2006-01-02", fechaInicioQ)
	if err != nil {
		fechaInicio = time.Date(1970, time.December, 28, 23, 59, 0, 0, tiempo.Local)
	}

	fechaFin, err := time.Parse("2006-01-02", fechaFinQ)
	if err != nil {
		fechaFin = time.Date(3000, time.December, 28, 23, 59, 0, 0, tiempo.Local)
	}
	aportaciones := []*models.Aportacion{}
	err = models.Db.Where("created_at between ? and ?", fechaInicio, fechaFin).Where(&models.Aportacion{DonacionID: uint(idD)}).Preload("Fiel").Preload("Fiel.Usuario").Find(&aportaciones).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener aportaciones"), nil, c, http.StatusInternalServerError)
		return
	}
	total := &DonacionesTotal{}
	for _, don := range aportaciones {
		total.Monto += don.Monto
	}
	total.Donaciones = aportaciones
	total.Donacion = donacion
	utils.CrearRespuesta(nil, total, c, http.StatusOK)
}
