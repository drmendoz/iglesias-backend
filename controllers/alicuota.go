package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	// "github.com/ahmetb/go-linq/v3"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/paymentez"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func filter(ss []*models.Alicuota, test func(models.Alicuota) bool) (ret []*models.Alicuota) {
	for _, s := range ss {
		if test(*s) {
			ret = append(ret, s)
		}
	}
	return
}

type ReporteAlicuota struct {
	Pendientes float64 `json:"pendientes"`
	Vencidas   float64 `json:"vencidas"`
}

func ReporteAlicuotas(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	dia := c.Query("dia")
	mes := c.Query("mes")
	año := c.Query("año")

	alicuotas := []*models.Alicuota{}
	var err error
	var fechaInicio time.Time
	var fechaFin time.Time
	if dia != "" {
		fechaInicio, _ = time.Parse("2006-01-02", dia)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener alicuotas en el formato de la fecha"), nil, c, http.StatusInternalServerError)
			return
		}
		fechaFin = fechaInicio.AddDate(0, 0, 1)
	} else if mes != "" {
		fechaInicio, _ = time.Parse("2006-01-02", mes)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener alicuotas en el formato de la fecha"), nil, c, http.StatusInternalServerError)
			return
		}
		fechaFin = fechaInicio.AddDate(0, 1, 0)
	} else if año != "" {
		fechaInicio, _ = time.Parse("2006-01-02", año)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener alicuotas en el formato de la fecha"), nil, c, http.StatusInternalServerError)
			return
		}
		fechaFin = fechaInicio.AddDate(1, 0, 0)
	}
	err = models.Db.Joins("Casa", "EtapaID = ?", idEtapa).Where("fecha_pago between ? and ?", fechaInicio, fechaFin, idEtapa).Order("created_at desc").Find(&alicuotas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuotas"), nil, c, http.StatusInternalServerError)
		return
	}
	reporte := &ReporteAlicuota{Pendientes: 0, Vencidas: 0}
	for _, alicuota := range alicuotas {
		if alicuota.Estado == "PENDIENTE" {
			reporte.Pendientes += alicuota.Valor
		} else if alicuota.Estado == "VENCIDO" {
			reporte.Vencidas += alicuota.Valor
		}
	}
	utils.CrearRespuesta(err, reporte, c, http.StatusOK)
}

type GrupoAlicuotas struct {
	Alicuotas     []*models.Alicuota `json:"alicuotas"`
	TotalPagadas  float64            `json:"total_pagadas"`
	TotalVencidas float64            `json:"total_vencidas"`
}
type RespuestaAlicuotas struct {
	Extraordinarias GrupoAlicuotas `json:"extraordinarias"`
	Saldos          GrupoAlicuotas `json:"saldos"`
	Comunes         interface{}    `json:"comunes"`
}

func GetAlicuotas(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	mz := c.Query("mz")
	villa := c.Query("villa")
	estado := c.Query("estado")

	alicuotas := []*models.Alicuota{}
	var err error
	if idEtapa != 0 {
		err = models.Db.Preload("Casa", models.Db.Where(&models.Casa{Manzana: mz, Villa: villa, EtapaID: idEtapa})).Where(&models.Alicuota{Estado: estado}).Order("created_at desc").Find(&alicuotas).Error
		test := func(ali models.Alicuota) bool { return ali.Casa != nil }
		alicuotas = filter(alicuotas, test)
	} else {
		err = models.Db.Where(&models.Alicuota{Estado: estado}).Order("created_at desc").Find(&alicuotas).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuotas"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, alicuotas, c, http.StatusOK)
}

func GetAlicuotaPorId(c *gin.Context) {
	alicuota := &models.Alicuota{}
	id := c.Param("id")
	err := models.Db.First(alicuota, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Alicuota no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuota"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, alicuota, c, http.StatusOK)
}

func CreateAlicuotaBulk(c *gin.Context) {
	alicuotas := []*models.Alicuota{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear alicuota"), nil, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &alicuotas)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear alicuota"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, alicuota := range alicuotas {
		err = models.Db.Create(alicuota).Error
		alicuota.MesPago = tiempo.BeginningOfMonth(*alicuota.FechaPago)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear alicuota"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	utils.CrearRespuesta(err, "Alicuotas creada correctamente", c, http.StatusCreated)
}

func CreateAlicuota(c *gin.Context) {
	alicuota := &models.Alicuota{}
	err := c.ShouldBindJSON(alicuota)
	alicuota.MesPago = tiempo.BeginningOfMonth(*alicuota.FechaPago)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Create(alicuota).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear alicuota"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Alicuota creada correctamente", c, http.StatusCreated)
}

func UpdateAlicuotaBulk(c *gin.Context) {
	alicuotas := []*models.Alicuota{}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear alicuota"), nil, c, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &alicuotas)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear alicuota"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, alicuota := range alicuotas {
		err = models.Db.Where("id = ?", alicuota.ID).Updates(alicuota).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar alicuota"), nil, c, http.StatusInternalServerError)
			return
		}
	}

	utils.CrearRespuesta(err, "Alicuotas actualizadas correctamente", c, http.StatusOK)
}

func UpdateAlicuota(c *gin.Context) {
	alicuota := &models.Alicuota{}

	err := c.ShouldBindJSON(alicuota)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(alicuota).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar alicuota"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Alicuota actualizada correctamente", c, http.StatusOK)
}

func DeleteAlicuota(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Alicuota{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar alicuota"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Alicuota eliminada exitosamente", c, http.StatusOK)
}

func GetAlicuotaPorCasa(c *gin.Context) {
	idCasa := c.Param("id")
	idCasaResidente := c.GetInt("id_casa")
	if idCasaResidente != 0 {
		idCasa = fmt.Sprintf("%d", idCasaResidente)
	}
	fechaInicio := c.Query("fecha_inicio")
	if fechaInicio == "" {
		fechaInicio = "1900-01-01"
	}
	fechaFin := c.Query("fecha_fin")
	if fechaFin == "" {
		fechaFin = "2500-01-01"
	}
	estado := c.Query("estado")
	alicuotas := []*models.Alicuota{}
	err := models.Db.Where("casa_id = ? and estado LIKE ? and created_at between ? and ?", idCasa, "%"+estado+"%", fechaInicio, fechaFin).Find(&alicuotas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuotas"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, alicuotas, c, http.StatusOK)

}

type AlicuotasResidenteTemporal struct {
	Estado         Estado         `json:"estado"`
	Alicuotas      []*AlicuotaAno `json:"alicuotas"`
	PagoHabilitado bool           `json:"pago_habilitado"`
	ExentoPago     bool           `json:"exento_pago"`
}

type Estado struct {
	Pendientes bool    `json:"pendientes"`
	Valor      float64 `json:"valor"`
}

type AlicuotaAno struct {
	Ano       int                `json:"ano"`
	Alicuotas []*models.Alicuota `json:"alicuotas"`
}

type AlicuotasResidente struct {
	Estado    Estado             `json:"estado"`
	Alicuotas []*models.Alicuota `json:"alicuotas"`
}

func GetAlicuotaPorResidente(c *gin.Context) {
	idCasaResidente := c.GetInt("id_casa")
	alicuotas := []*models.Alicuota{}
	err := models.Db.Order("fecha_pago ASC").Where("casa_id = ? ", idCasaResidente).Find(&alicuotas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuotas"), nil, c, http.StatusInternalServerError)
		return
	}
	aliResidente := &AlicuotasResidente{}
	aliResidente.Alicuotas = alicuotas
	for _, ali := range alicuotas {
		if ali.Estado == "PENDIENTE" {
			aliResidente.Estado.Valor += ali.Valor
		}
	}

	aliResidente.Estado.Valor = math.Round(aliResidente.Estado.Valor*100) / 100
	if aliResidente.Estado.Valor > 0 {
		aliResidente.Estado.Pendientes = true
	}
	utils.CrearRespuesta(err, aliResidente, c, http.StatusOK)

}

func GetAlicuotaPorResidenteTemporal(c *gin.Context) {
	idCasaResidente := c.GetInt("id_casa")
	alicuotas := []*models.Alicuota{}
	err := models.Db.Order("fecha_pago ASC").Where("casa_id = ? ", idCasaResidente).Find(&alicuotas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuotas"), nil, c, http.StatusInternalServerError)
		return
	}
	idEtapa := c.GetInt("id_etapa")
	etapa := &models.Etapa{}
	err = models.Db.Find(&etapa, idEtapa).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener alicuotas"), nil, c, http.StatusInternalServerError)
		return
	}
	aliResidente := &AlicuotasResidenteTemporal{}
	aliResidente.Alicuotas = []*AlicuotaAno{}
	vencidos := false
	for _, ali := range alicuotas {
		if ali.Estado == "PENDIENTE" {
			aliResidente.Estado.Valor += ali.Valor
			if time.Now().After(*ali.FechaPago) {
				ali.Estado = "VENCIDO"
				vencidos = true
			}
		}

		ano := ali.FechaPago.Year()
		flag := true
		for _, aliAno := range aliResidente.Alicuotas {
			if aliAno.Ano == ano {
				aliAno.Alicuotas = append(aliAno.Alicuotas, ali)
				flag = false
				break
			}
		}
		aliResidente.Estado.Valor = math.Round(aliResidente.Estado.Valor*100) / 100
		if flag {
			alIt := []*models.Alicuota{}
			alIt = append(alIt, ali)
			aliResidente.Alicuotas = append(aliResidente.Alicuotas, &AlicuotaAno{Ano: ano, Alicuotas: alIt})
		}
	}
	aliResidente.ExentoPago = c.GetInt("id_residente") == 34
	aliResidente.Estado.Pendientes = vencidos
	aliResidente.PagoHabilitado = etapa.PagosTarjeta
	utils.CrearRespuesta(err, aliResidente, c, http.StatusOK)

}

func CreateAlicuotaAutomaticamente(c *gin.Context) {
	etapas := []*models.Etapa{}
	err := models.Db.Preload("Casa").Find(&etapas).Error
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al crear alicoutas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, etapa := range etapas {
		utils.Log.Warn(etapa.ValorAlicuota)
		for _, casa := range etapa.Casa {
			date := time.Now().In(tiempo.Local)
			for i := int(date.Month()); i < 13; i++ {
				//Se iguala a 1 para que se genere el cobro de alicuota cada 1 del mes por solicitud del cliente
				etapa.FechaAlicuota = 1
				datePagoInicio := time.Date(date.Year(), date.Month(), etapa.FechaAlicuota, 0, 0, 0, 0, tiempo.Local)
				datePagoFinal := time.Date(date.Year(), date.Month(), etapa.FechaAlicuota, 23, 59, 59, 59, tiempo.Local)
				alicuota := &models.Alicuota{}
				err = models.Db.Where("fecha_pago between ? and  ? and casa_id = ?", datePagoInicio, datePagoFinal, casa.ID).First(alicuota).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						err = models.Db.Create(&models.Alicuota{Valor: etapa.ValorAlicuota, CasaID: casa.ID, FechaPago: &datePagoInicio}).Error
						if err != nil {
							utils.CrearRespuesta(errors.New("Error al crear alicuotas"), nil, c, http.StatusInternalServerError)
							return
						}
					}
				}
				date = date.AddDate(0, 1, 0)
			}
		}
	}

	utils.CrearRespuesta(nil, "Alicuotas creadas exitosamente", c, http.StatusOK)
}

type PagoAlicuota struct {
	IdAlicuotas  []uint    `json:"alicuotas"`
	TokenTarjeta string    `json:"token_tarjeta"`
	FechaPago    time.Time `json:"fecha_pago"`
}

func PagarAlicuota(c *gin.Context) {
	pago := &PagoAlicuota{}
	err := c.ShouldBindJSON(pago)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en formato de respuesta"), nil, c, http.StatusInternalServerError)
		return
	}
	pago.FechaPago = time.Now()
	tx := models.Db.Begin()
	idResidente := c.GetInt("id_residente")
	idRes := fmt.Sprintf("%d", idResidente)
	res := &models.Residente{}
	err = tx.Joins("Usuario").Find(res, idRes).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	descripcion := ""
	monto := 0.0
	for _, idAlicuota := range pago.IdAlicuotas {
		ali := &models.Alicuota{}
		err = tx.Where("id = ?", idAlicuota).First(ali).Error
		if err != nil {
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		}
		if ali.Estado == "PAGADO" {
			respuesta := fmt.Sprintf("La alicuota con identificador %d ya ha sido pagada", idAlicuota)
			tx.Rollback()
			utils.CrearRespuesta(errors.New(respuesta), nil, c, http.StatusNotAcceptable)
			return
		}
		err = tx.Where("id = ?", idAlicuota).Updates(&models.Alicuota{Estado: "Pagado"}).Error
		descripcion += fmt.Sprintf("Pago de alicuota id # %d. ", idAlicuota)
		monto += ali.Valor
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar alicuotas"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tarjeta := &models.ResidenteTarjeta{}
	err = tx.Where("token_tarjeta = ?", pago.TokenTarjeta).First(tarjeta).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tarjeta.TokenTarjeta = pago.TokenTarjeta
			tarjeta.ResidenteID = res.ID
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

	trans := &models.Transaccion{ResidenteTarjetaID: tarjeta.ID, Tipo: "ALI"}
	err = tx.Create(trans).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	idTrans := fmt.Sprintf("%d", trans.ID)
	cobro, err := paymentez.CobrarTarjeta(idRes, res.Usuario.Correo, monto, descripcion, idTrans, 0, pago.TokenTarjeta)
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al debitar tarjeta"), nil, c, http.StatusPaymentRequired)
		return
	}
	montoReal := fmt.Sprintf("%f", cobro.Transaccion.Monto)
	transNueva := &models.Transaccion{Estado: cobro.Transaccion.Status, DiaPago: cobro.Transaccion.FechaPago, Monto: montoReal, CodigoAutorizacion: cobro.Transaccion.CodigoAutorizacion, Mensaje: cobro.Transaccion.Mensaje, Descripcion: descripcion, ResidenteTarjetaID: tarjeta.ID}
	err = tx.Where("id = ?", trans.ID).Updates(transNueva).Error
	if err != nil {
		_ = c.Error(err)
		_ = c.Error(errors.New("Alicuotas pagadas pero error al guardar transaccion"))
		tx.Commit()
		utils.CrearRespuesta(nil, "Alicuotas pagadas exitosamente", c, http.StatusOK)
		return
	}
	for _, idAlicuota := range pago.IdAlicuotas {
		err = tx.Where("id= ?", idAlicuota).Updates(&models.Alicuota{TransaccionID: &trans.ID}).Error
		if err != nil {
			_ = c.Error(err)
		}
	}

	tx.Commit()
	utils.CrearRespuesta(nil, "Alicuotas pagadas exitosamente", c, http.StatusOK)
}

func DevolverPagoAlicuotas(c *gin.Context) {}
