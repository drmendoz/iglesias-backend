package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetParroquias(c *gin.Context) {
	etps := []*models.Parroquia{}
	err := models.Db.Order("Nombre ASC").Preload("Iglesia").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener parroquias"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, etp := range etps {
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultParroquia
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}

	}
	utils.CrearRespuesta(err, etps, c, http.StatusOK)
}

func GetParroquiaPorId(c *gin.Context) {
	etp := &models.Parroquia{}
	id := c.Param("id")
	err := models.Db.First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Parroquia no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultParroquia
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateParroquia(c *gin.Context) {
	etp := &models.Parroquia{}
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
		etp.Imagen = utils.DefaultParroquia
	} else {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "parroquias/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(etp.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Parroquia{}).Where("id = ?", etp.ID).Update("imagen", etp.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)
			return
		}
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Parroquia creada correctamente", c, http.StatusCreated)

}

func UpdateParroquia(c *gin.Context) {
	etp := &models.Parroquia{}

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
	// err = tx.Model(etp).Where("id = ?", id).Updates(map[string]interface{}{
	// 	"pagos_tarjeta":         etp.PagosTarjeta,
	// 	"modulo_market":         etp.ModuloMarket,
	// 	"modulo_publicacion":    etp.ModuloPublicacion,
	// 	"modulo_votacion":       etp.ModuloVotacion,
	// 	"modulo_area_social":    etp.ModuloAreaSocial,
	// 	"modulo_equipo":         etp.ModuloEquipoAdministrativo,
	// 	"modulo_historia":       etp.ModuloHistoria,
	// 	"modulo_bitacora":       etp.ModuloBitacora,
	// 	"formulario_entrada":    etp.FormularioEntrada,
	// 	"formulario_salida":     etp.FormularioSalida,
	// 	"modulo_alicuota":       etp.ModuloAlicuota,
	// 	"modulo_emprendimiento": etp.ModuloEmprendimiento,
	// 	"modulo_camaras":        etp.ModuloCamaras,
	// 	"modulo_directiva":      etp.ModuloDirectiva,
	// 	"modulo_galeria":        etp.ModuloGaleria,
	// 	"modulo_horarios":       etp.ModuloHorarios,
	// 	"modulo_mi_registro":    etp.ModuloMiRegistro}).Error

	// if err != nil {
	// 	_ = c.Error(err)
	// 	tx.Rollback()
	// 	utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
	// 	return
	// }
	if etp.Imagen != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "parroquias/"+time.RFC3339+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear etapa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Parroquia{}).Where("id = ?", etp.ID).Update("imagen", etp.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al actualizar etapa"), nil, c, http.StatusInternalServerError)
			return
		}
		etp.Imagen = utils.SERVIMG + etp.Imagen

	} else {
		etp.Imagen = utils.DefaultParroquia
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Parroquia actualizada correctamente", c, http.StatusOK)
}

func DeleteParroquia(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Parroquia{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar etapa"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Parroquia eliminada exitosamente", c, http.StatusOK)
}
