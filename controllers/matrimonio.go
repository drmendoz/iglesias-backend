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

func GetMatrimonios(c *gin.Context) {
	etps := []*models.Matrimonio{}
	estado := c.Query("estado")
	idParroquia := c.GetInt("id_parroquia")
	err := models.Db.Where(&models.Matrimonio{ParroquiaID: uint(idParroquia), Estado: estado}).Order("created_at asc").Find(&etps).Error

	for _, mat := range etps {
		mat.Imagen = utils.SERVIMG + mat.Imagen

	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener matrimonios"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, etps, c, http.StatusOK)
}

func GetMatrimonioPorID(c *gin.Context) {
	mat := &models.Matrimonio{}
	id := c.Param("id")
	err := models.Db.Preload("MatrimonioImagenes").First(mat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Matrimonio no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener curso"), nil, c, http.StatusInternalServerError)
		return
	}
	mat.Imagen = utils.SERVIMG + mat.Imagen

	utils.CrearRespuesta(nil, mat, c, http.StatusOK)
}

func CreateMatrimonio(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_parroquia"))
	idFiel := uint(c.GetInt("id_fiel"))
	mat := &models.Matrimonio{}
	err := c.ShouldBindJSON(mat)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	mat.ParroquiaID = idParroquia
	mat.FielID = idFiel

	tx := models.Db.Begin()
	//if mat.Monto == 0 {
	mat.Estado = "PEN"
	err = tx.Omit("imagen").Create(mat).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear imagen"), nil, c, http.StatusInternalServerError)
		return
	}

	if mat.Imagen == "" {
		mat.Imagen = utils.DefaultGaleria
	} else {
		idUrb := fmt.Sprintf("%d", mat.ID)
		mat.Imagen, err = img.FromBase64ToImage(mat.Imagen, "imagens/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear imagen "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Matrimonio{}).Where("id = ?", mat.ID).Update("imagen", mat.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear imagen "), nil, c, http.StatusInternalServerError)
			return
		}
		mat.Imagen = utils.SERVIMG + mat.Imagen
	}
	_ = tx.Commit()
	// if mat.TokenTarjeta == "" {
	// 	_ = tx.Rollback()
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Ingrese tarjeta de credito"), nil, c, http.StatusNotAcceptable)
	// 	return
	// }
	// fiel := &models.Fiel{}
	// err = tx.Joins("Usuario").Find(fiel, idFiel).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
	// 	return
	// }
	// tarjeta := &models.FielTarjeta{TokenTarjeta: mat.TokenTarjeta, FielID: uint(idFiel)}
	// err = tx.FirstOrCreate(tarjeta).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
	// 	return
	// }
	// mat.Transaccion = &models.Transaccion{FielTarjetaID: tarjeta.ID, ParroquiaID: idParroquia}
	// err = tx.Create(mat).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
	// 	return
	// }
	// idTrans := fmt.Sprintf("%d", mat.Transaccion.ID)
	// descripcion := fmt.Sprintf("Pago de Matrimonio # %d", mat.ID)
	// idFi := fmt.Sprintf("%d", idFiel)
	// cobro, err := paymentez.CobrarTarjeta(idFi, fiel.Usuario.Correo, mat.Monto, descripcion, idTrans, 0, mat.TokenTarjeta)
	// if err != nil {
	// 	tx.Rollback()
	// 	_ = c.Error(err)
	// 	utils.CrearRespuesta(errors.New("Error al debitar tarjeta"), nil, c, http.StatusPaymentRequired)
	// 	return
	// }
	// montoReal := fmt.Sprintf("%f", cobro.Transaccion.Monto)
	// transNueva := &models.Transaccion{Estado: cobro.Transaccion.Status, DiaPago: cobro.Transaccion.FechaPago, Monto: montoReal, CodigoAutorizacion: cobro.Transaccion.CodigoAutorizacion, Mensaje: cobro.Transaccion.Mensaje, Descripcion: descripcion, FielTarjetaID: tarjeta.ID}
	// err = tx.Where("id = ?", mat.Transaccion.ID).Updates(transNueva).Error
	// if err != nil {
	// 	_ = c.Error(err)
	// 	tx.Rollback()
	// 	utils.CrearRespuesta(errors.New("Error al registrar matrimonio"), nil, c, http.StatusOK)
	// 	return
	// }
	// _ = tx.Commit()
	utils.CrearRespuesta(nil, "Matrimonio subido exitosamente", c, http.StatusOK)

}

func UpdateMatrimonio(c *gin.Context) {
	mat := &models.Matrimonio{}

	err := c.ShouldBindJSON(mat)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	id := c.Param("id")
	err = tx.Where("id = ?", id).Updates(mat).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar matrimonio"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Matrimonio actualizado correctamente", c, http.StatusOK)
}

func DeleteMatrimonio(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Matrimonio{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar curso"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Matrimonio eliminada exitosamente", c, http.StatusOK)
}

func GetMatrimoniosUsuario(c *gin.Context) {

	etps := []*models.Matrimonio{}
	idFiel := uint(c.GetInt("id_fiel"))
	idParroquia := uint(c.GetInt("id_parroquia"))
	err := models.Db.Where(&models.Matrimonio{FielID: idFiel, ParroquiaID: idParroquia}).Preload("MatrimonioImagenes").Find(&etps).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Matrimonio no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener curso"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, mat := range etps {
		mat.Imagen = utils.SERVIMG + mat.Imagen
	}
	utils.CrearRespuesta(nil, etps, c, http.StatusOK)
}
