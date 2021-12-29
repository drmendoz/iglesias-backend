package controllers

import (
	"errors"
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

func GetTransaccions(c *gin.Context) {
	modulo := c.Query("modulo")
	idCategoria := c.Query("id_categoria")
	idParroquia := c.Query("id_parroquia")
	idCaso := c.Query("id_caso")
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
	idCat, err := strconv.Atoi(idCategoria)
	if err != nil {
		idCat = 0
	}

	idPar, err := strconv.Atoi(idParroquia)
	if err != nil {
		idPar = 0
	}
	idCas, err := strconv.Atoi(idCaso)
	if err != nil {
		idCas = 0
	}
	idParr := c.GetInt("id_parroquia")
	if idParr != 0 {
		idPar = idParr
	}
	transaccions := []*models.Transaccion{}
	err = models.Db.Where("created_at between ? and ?", fechaInicio, fechaFin).Where(&models.Transaccion{TipoPagoType: modulo, ParroquiaID: uint(idPar), CategoriaID: uint(idCat), CasoID: uint(idCas)}).Preload("FielTarjeta.Fiel.Usuario").Find(&transaccions).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener transacciones"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, tr := range transaccions {
		if tr.FielTarjeta.Fiel != nil {
			tr.Nombre = tr.FielTarjeta.Fiel.Usuario.Nombre
			tr.Apellido = tr.FielTarjeta.Fiel.Usuario.Apellido
			tr.Correo = tr.FielTarjeta.Fiel.Usuario.Correo
		}
		tr.FielTarjeta = nil
		if tr.CategoriaID != 0 {
			if tr.TipoPagoType == "aportacion" {
				cat := &models.CategoriaDonacion{}
				err = models.Db.First(cat, tr.CategoriaID).Error
				if err != nil {
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error al obtener transacciones"), nil, c, http.StatusInternalServerError)
					return
				}
				tr.NombreCategoria = cat.Nombre
			} else if tr.TipoPagoType == "emprendimiento" {
				cat := &models.CategoriaMarket{}
				err = models.Db.First(cat, tr.CategoriaID).Error
				if err != nil {
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error al obtener transacciones"), nil, c, http.StatusInternalServerError)
					return
				}
				tr.NombreCategoria = cat.Nombre
			}
		}
	}
	utils.CrearRespuesta(nil, transaccions, c, http.StatusOK)
}

func GetTransaccion(c *gin.Context) {
	transaccions := []*models.Transaccion{}
	err := models.Db.Order("created_at DESC").Find(&transaccions).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener transacciones"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, transaccions, c, http.StatusOK)
}

func GetTransaccionPorId(c *gin.Context) {
	transaccion := &models.Transaccion{}
	id := c.Param("id")
	err := models.Db.First(transaccion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Transaccion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener transaccion"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, transaccion, c, http.StatusOK)
}

func DevolverTransaccion(c *gin.Context) {
	id := c.Param("id")
	tx := models.Db.Begin()
	transaccion := &models.Transaccion{}
	err := models.Db.First(transaccion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			utils.CrearRespuesta(errors.New("No existe transaccion"), nil, c, http.StatusBadRequest)
			return
		}
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al devolver transaccion"), nil, c, http.StatusInternalServerError)
		return
	}
	respuesta, err := paymentez.DevolverPago(transaccion.ID)
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Servicio no disponible"), nil, c, http.StatusInternalServerError)
		return
	}
	estado := ""
	if respuesta.Status == "success" {
		estado = "Devuelta"
	} else if respuesta.Status == "pending" {
		estado = " Devolucion pendiente"
	}
	err = tx.Where("id = ?", transaccion.ID).Updates(&models.Transaccion{Estado: "Devuelto", EstadoDevolucion: estado, DetalleDevolucion: respuesta.Detalle}).Error
	if err != nil {
		_ = c.Error(errors.New("Devolucion entrega, pero error al cambiar estado de transaccion"))

	}
	tx.Commit()
	utils.CrearRespuesta(nil, "Transaccion devuelta", c, http.StatusOK)

}

func CreateTransaccion(c *gin.Context) {
	transaccion := &models.Transaccion{}
	err := c.ShouldBindJSON(transaccion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Create(transaccion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear transaccion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, "Transaccion creada con exito", c, http.StatusCreated)

}

func UpdateTransaccion(c *gin.Context) {
	transaccion := &models.Transaccion{}

	err := c.ShouldBindJSON(transaccion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	err = models.Db.Where("id = ?", id).Updates(transaccion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar transaccion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, "Transaccion actualizada correctamente", c, http.StatusOK)
}

func DeleteTransaccion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Transaccion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar transaccion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Transaccion eliminada exitosamente", c, http.StatusOK)
}
