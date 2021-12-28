package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/paymentez"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetIntenciones(c *gin.Context) {
	idFiel := c.GetInt("id_fiel")
	idParroquia := c.GetInt("id_parroquia")
	idMisa := c.Query("id_misa")
	idMis, err := strconv.Atoi(idMisa)
	if err != nil {
		idMis = 0
	}
	etps := []*models.Intencion{}
	err = models.Db.Where(&models.Intencion{ParroquiaID: uint(idParroquia), FielID: uint(idFiel), MisaID: uint(idMis)}).Order("created_at ASC").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener intencions"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetIntencionPorId(c *gin.Context) {
	etp := &models.Intencion{}
	id := c.Param("id")
	err := models.Db.First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Intención no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener intención"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateIntencion(c *gin.Context) {
	etp := &models.Intencion{}
	idParroquia := c.GetInt("id_parroquia")
	idFiel := c.GetInt("id_fiel")

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = uint(idParroquia)
	etp.FielID = uint(idFiel)

	tx := models.Db.Begin()
	misa := &models.Misa{}
	err = tx.First(misa, etp.MisaID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = tx.Rollback()
			utils.CrearRespuesta(errors.New("No existe misa"), nil, c, http.StatusBadRequest)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear intencion"), nil, c, http.StatusInternalServerError)
		return

	}
	var numIntencionesMisa int64
	err = tx.Model(&models.Intencion{}).Where(&models.Intencion{MisaID: etp.MisaID}).Count(&numIntencionesMisa).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear intención"), nil, c, http.StatusInternalServerError)
		return
	}
	if numIntencionesMisa >= int64(misa.CupoIntencion) {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Limite de intenciones cumplidas"), nil, c, http.StatusNotAcceptable)
		return
	}
	parroquia := &models.Parroquia{}
	err = tx.First(parroquia, idParroquia).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear intención"), nil, c, http.StatusInternalServerError)
		return
	}
	if !parroquia.BotonPagoIntencion || etp.Monto == 0 {
		err = tx.Create(etp).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear intención"), nil, c, http.StatusInternalServerError)
			return
		}

		tx.Commit()
		utils.CrearRespuesta(err, "Intención creada correctamente", c, http.StatusCreated)
		return
	}

	if etp.TokenTarjeta == "" {
		_ = tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Ingrese tarjeta de credito"), nil, c, http.StatusNotAcceptable)
		return
	}
	fiel := &models.Fiel{}
	err = tx.Joins("Usuario").Find(fiel, idFiel).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	tarjeta := &models.FielTarjeta{TokenTarjeta: etp.TokenTarjeta, FielID: uint(idFiel)}
	err = tx.FirstOrCreate(tarjeta).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear intencion"), nil, c, http.StatusInternalServerError)
		return
	}
	etp.Transaccion = &models.Transaccion{FielTarjetaID: tarjeta.ID, ParroquiaID: uint(idParroquia)}
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al cerear intecion"), nil, c, http.StatusInternalServerError)
		return
	}
	idTrans := fmt.Sprintf("%d", etp.Transaccion.ID)
	descripcion := fmt.Sprintf("Pago de intecion # %d", etp.ID)
	idFi := fmt.Sprintf("%d", idFiel)
	cobro, err := paymentez.CobrarTarjeta(idFi, fiel.Usuario.Correo, etp.Monto, descripcion, idTrans, 0, etp.TokenTarjeta)
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al debitar tarjeta"), nil, c, http.StatusPaymentRequired)
		return
	}
	montoReal := fmt.Sprintf("%f", cobro.Transaccion.Monto)
	transNueva := &models.Transaccion{Estado: cobro.Transaccion.Status, DiaPago: cobro.Transaccion.FechaPago, Monto: montoReal, CodigoAutorizacion: cobro.Transaccion.CodigoAutorizacion, Mensaje: cobro.Transaccion.Mensaje, Descripcion: descripcion, FielTarjetaID: tarjeta.ID}
	err = tx.Where("id = ?", etp.Transaccion.ID).Updates(transNueva).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al crear intecion"), nil, c, http.StatusOK)
		return
	}
	_ = tx.Commit()
	utils.CrearRespuesta(nil, "Intecion subido exitosamente", c, http.StatusOK)

}

func UpdateIntencion(c *gin.Context) {
	etp := &models.Intencion{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar intención"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Intención actualizada correctamente", c, http.StatusOK)
}

func DeleteIntencion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Intencion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar intención"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Intención eliminada exitosamente", c, http.StatusOK)
}
