package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCasas(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	mz := c.Query("mz")

	casas := []*models.Casa{}
	err := models.Db.Order("manzana asc").Order("villa asc").Where(&models.Casa{EtapaID: idEtapa, Manzana: mz}).Preload("Alicuotas").Find(&casas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener casas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, casa := range casas {
		casa.Etapa = nil
		if casa.Imagen == "" {
			casa.Imagen = utils.DefaultCasa
		} else {
			casa.Imagen = utils.SERVIMG + casa.Imagen
		}
		casa.DebeAlicuotas = false
		casa.ValorAlicuotas = 0
		for _, ali := range casa.Alicuotas {
			if ali.Estado == "VENCIDO" && ali.Valor > 0 {
				casa.DebeAlicuotas = true
				casa.ValorAlicuotas = casa.ValorAlicuotas + ali.Valor
			}
		}
		casa.Alicuotas = nil
	}

	utils.CrearRespuesta(err, casas, c, http.StatusOK)
}

func GetCasaPorId(c *gin.Context) {
	casa := &models.Casa{}
	id := c.Param("id")
	err := models.Db.Preload("Alicuotas").First(casa, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Casa no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener casa"), nil, c, http.StatusInternalServerError)
		return
	}
	if casa.Imagen == "" {
		casa.Imagen = utils.DefaultCasa
	} else {
		casa.Imagen = utils.SERVIMG + casa.Imagen
	}
	casa.DebeAlicuotas = false
	casa.ValorAlicuotas = 0
	for _, ali := range casa.Alicuotas {
		if ali.Estado == "PENDIENTE" && ali.Valor > 0 {
			casa.DebeAlicuotas = true
			casa.ValorAlicuotas = casa.ValorAlicuotas + ali.Valor
		}
	}

	utils.CrearRespuesta(nil, casa, c, http.StatusOK)
}

func CreateCasa(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	casa := &models.Casa{}
	err := c.ShouldBindJSON(casa)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	casa.EtapaID = idEtapa
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(casa).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear casa"), nil, c, http.StatusInternalServerError)
		return
	}

	if casa.Imagen == "" {
		casa.Imagen = utils.DefaultCasa
	} else {
		idUrb := fmt.Sprintf("%d", casa.ID)
		casa.Imagen, err = img.FromBase64ToImage(casa.Imagen, "casas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(casa.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear casa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Casa{}).Where("id = ?", casa.ID).Update("imagen", casa.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear casa "), nil, c, http.StatusInternalServerError)
			return
		}
		casa.Imagen = utils.SERVIMG + casa.Imagen
	}
	err = createAlicuotasCasa(tx, casa)
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al crear casa"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Casa creada correctamente", c, http.StatusCreated)

}

func createAlicuotasCasa(tx *gorm.DB, casa *models.Casa) error {
	etapa := &models.Etapa{}
	err := models.Db.First(etapa, casa.EtapaID).Error
	date := time.Now().In(tiempo.Local)
	mes := int(date.Month())
	for i := mes; i < 13; i++ {
		datePagoInicio := time.Date(date.Year(), date.Month(), etapa.FechaAlicuota, 0, 0, 0, 0, tiempo.Local)
		err = tx.Create(&models.Alicuota{Valor: etapa.ValorAlicuota, CasaID: casa.ID, FechaPago: &datePagoInicio}).Error
		if err != nil {
			return err
		}
		date = date.AddDate(0, 1, 0)
	}
	return nil
}

func UpdateCasa(c *gin.Context) {
	casa := &models.Casa{}

	err := c.ShouldBindJSON(casa)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("imagen", "id_etapa").Where("id = ?", id).Updates(casa).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar casa"), nil, c, http.StatusInternalServerError)
		return
	}
	if casa.Imagen != "" {
		idUrb := fmt.Sprintf("%d", casa.ID)
		casa.Imagen, err = img.FromBase64ToImage(casa.Imagen, "casas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear casa "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Casa{}).Where("id = ?", casa.ID).Update("imagen", casa.Imagen).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar casa"), nil, c, http.StatusInternalServerError)
			return
		}
		casa.Imagen = utils.SERVIMG + casa.Imagen

	} else {
		casa.Imagen = utils.DefaultCasa
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Casa actualizada correctamente", c, http.StatusOK)
}

func DeleteCasa(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Casa{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar casa"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Casa eliminada exitosamente", c, http.StatusOK)
}

func GetCasasPorEtapa(c *gin.Context) {
	idEtapa := c.Param("id")
	casas := []*models.Casa{}
	err := models.Db.Where("etapa_id = ?", idEtapa).Find(&casas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener casas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, casa := range casas {
		if casa.Imagen == "" {
			casa.Imagen = utils.SERVIMG + "default_user.png"
		} else {
			casa.Imagen = utils.SERVIMG + casa.Imagen
		}

	}
	utils.CrearRespuesta(err, casas, c, http.StatusOK)
}

func GetCasasPorUrbanizacion(c *gin.Context) {
	idUrbanizacion := c.Param("id")
	casas := []*models.Casa{}
	err := models.Db.Where("Etapa.urbanizacion_id = ?", idUrbanizacion).Joins("Etapa").Find(&casas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener casas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, casa := range casas {
		if casa.Imagen == "" {
			casa.Imagen = utils.SERVIMG + "default_user.png"
		} else {
			casa.Imagen = utils.SERVIMG + casa.Imagen
		}

	}
	utils.CrearRespuesta(err, casas, c, http.StatusOK)
}

func GetCasasCount(c *gin.Context) {
	var casas int64
	err := models.Db.Model(&models.Casa{}).Count(&casas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, casas, c, http.StatusOK)
}
