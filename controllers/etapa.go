package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEtapas(c *gin.Context) {
	urbanizacion, err := strconv.Atoi(c.Query("id_urbanizacion"))
	etps := []*models.Etapa{}
	err = models.Db.Where(&models.Etapa{UrbanizacionID: uint(urbanizacion)}).Joins("Urbanizacion").Order("Nombre ASC").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener etapas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, etp := range etps {
		etp.NombreUrbanizacion = etp.Urbanizacion.Nombre
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultEtapa
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}
		etp.Urbanizacion = nil

	}
	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetEtapaPorId(c *gin.Context) {
	etp := &models.Etapa{}
	id := c.Param("id")
	err := models.Db.First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Etapa no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultEtapa
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateEtapa(c *gin.Context) {
	etp := &models.Etapa{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear etapa"), nil, c, http.StatusInternalServerError)
		return
	}

	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultEtapa
	} else {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "etapas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(etp.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Etapa{}).Where("id = ?", etp.ID).Update("imagen", etp.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)
			return
		}
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Etapa creada correctamente", c, http.StatusCreated)

}

func UpdateEtapa(c *gin.Context) {
	etp := &models.Etapa{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(etp).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	err = tx.Model(etp).Where("id = ?", id).Updates(map[string]interface{}{
		"pagos_tarjeta":         etp.PagosTarjeta,
		"modulo_market":         etp.ModuloMarket,
		"modulo_publicacion":    etp.ModuloPublicacion,
		"modulo_votacion":       etp.ModuloVotacion,
		"modulo_area_social":    etp.ModuloAreaSocial,
		"modulo_equipo":         etp.ModuloEquipoAdministrativo,
		"modulo_historia":       etp.ModuloHistoria,
		"modulo_bitacora":       etp.ModuloBitacora,
		"formulario_entrada":    etp.FormularioEntrada,
		"formulario_salida":     etp.FormularioSalida,
		"modulo_alicuota":       etp.ModuloAlicuota,
		"modulo_emprendimiento": etp.ModuloEmprendimiento,
		"modulo_camaras":        etp.ModuloCamaras,
		"modulo_directiva":      etp.ModuloDirectiva,
		"modulo_galeria":        etp.ModuloGaleria,
		"modulo_horarios":       etp.ModuloHorarios,
		"modulo_mi_registro":    etp.ModuloMiRegistro}).Error

	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "etapas/"+time.RFC3339+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Etapa{}).Where("id = ?", etp.ID).Update("imagen", etp.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
			return
		}
		etp.Imagen = utils.SERVIMG + etp.Imagen

	} else {
		etp.Imagen = utils.DefaultEtapa
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Etapa actualizada correctamente", c, http.StatusOK)
}

func DeleteEtapa(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Etapa{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Etapa eliminada exitosamente", c, http.StatusOK)
}
