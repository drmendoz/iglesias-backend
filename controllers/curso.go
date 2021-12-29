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

func GetCursos(c *gin.Context) {
	etps := []*models.Curso{}
	idParroquia := c.GetInt("id_parroquia")
	err := models.Db.Where(&models.Curso{ParroquiaID: uint(idParroquia)}).Preload("Inscritos").Order("created_at asc").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener areas sociales"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, curs := range etps {
		curs.Imagen = utils.SERVIMG + curs.Imagen
		curs.Video = utils.SERVIMG + curs.Video
	}
	utils.CrearRespuesta(nil, etps, c, http.StatusOK)
}

func GetCursoPorID(c *gin.Context) {
	etp := &models.Curso{}
	id := c.Param("id")
	err := models.Db.Preload("Inscritos").First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Curso no encontrado"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener curso"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultCurso
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}

	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateCurso(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_parroquia"))
	etp := &models.Curso{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = idParroquia
	err = models.Db.Create(etp).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear curso"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Curso creada correctamente", c, http.StatusCreated)

}

func UpdateCurso(c *gin.Context) {
	etp := &models.Curso{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	id := c.Param("id")
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar curso"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Curso actualizada correctamente", c, http.StatusOK)
}

func DeleteCurso(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Curso{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar curso"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Curso eliminada exitosamente", c, http.StatusOK)
}

func InscribirCurso(c *gin.Context) {
	idFiel := c.GetInt("id_fiel")
	idParroquia := c.GetInt("id_parroquia")
	id := c.Param("id")
	idCurso, err := strconv.Atoi(id)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al inscribir a curso"), nil, c, http.StatusInternalServerError)
		return
	}
	inscrito := &models.Inscrito{}
	err = c.ShouldBindJSON(inscrito)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al inscribir a curso"), nil, c, http.StatusInternalServerError)
		return
	}

	tx := models.Db.Begin()
	curso := &models.Curso{}
	err = tx.First(curso, idCurso).Error
	if err != nil {
		_ = tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al inscribir a curso"), nil, c, http.StatusInternalServerError)
		return
	}
	inscrito.FielID = uint(idFiel)
	inscrito.CursoID = uint(idCurso)
	if curso.TieneLimite {
		var cupos int64
		err = tx.Model(&models.Inscrito{}).Where(&models.Inscrito{CursoID: uint(idCurso)}).Count(&cupos).Error
		if err != nil {
			_ = tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al inscribir a curso"), nil, c, http.StatusInternalServerError)
			return
		}
		if cupos >= int64(curso.Cupo) {
			_ = tx.Rollback()
			utils.CrearRespuesta(errors.New("Ya no hay cupo para este curso"), nil, c, http.StatusNotAcceptable)
		}
	}
	if !curso.BotonPago {
		err = tx.Create(inscrito).Error
		if err != nil {
			_ = tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al inscribir a curso"), nil, c, http.StatusInternalServerError)
			return
		}
		_ = tx.Commit()
		utils.CrearRespuesta(nil, "Inscripcion exitosa", c, http.StatusOK)
		return
	}

	if inscrito.TokenTarjeta == "" {
		_ = tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Ingrese tarjeta de credito"), nil, c, http.StatusNotAcceptable)
		return
	}
	inscrito.Monto = curso.Precio
	fiel := &models.Fiel{}
	err = tx.Joins("Usuario").Find(fiel, idFiel).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	tarjeta := &models.FielTarjeta{TokenTarjeta: inscrito.TokenTarjeta, FielID: uint(idFiel)}
	err = tx.FirstOrCreate(tarjeta).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	inscrito.Transaccion = &models.Transaccion{FielTarjetaID: tarjeta.ID, ParroquiaID: uint(idParroquia)}
	err = tx.Create(inscrito).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener informacion"), nil, c, http.StatusInternalServerError)
		return
	}
	idTrans := fmt.Sprintf("%d", inscrito.Transaccion.ID)
	descripcion := fmt.Sprintf("Pago  de curso: Inscripcion # %d", inscrito.ID)
	idFi := fmt.Sprintf("%d", idFiel)
	cobro, err := paymentez.CobrarTarjeta(idFi, fiel.Usuario.Correo, inscrito.Monto, descripcion, idTrans, 0, inscrito.TokenTarjeta)
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al debitar tarjeta"), nil, c, http.StatusPaymentRequired)
		return
	}
	montoReal := fmt.Sprintf("%f", cobro.Transaccion.Monto)
	transNueva := &models.Transaccion{Estado: cobro.Transaccion.Status, DiaPago: cobro.Transaccion.FechaPago, Monto: montoReal, CodigoAutorizacion: cobro.Transaccion.CodigoAutorizacion, Mensaje: cobro.Transaccion.Mensaje, Descripcion: descripcion, FielTarjetaID: tarjeta.ID}
	err = tx.Where("id = ?", inscrito.Transaccion.ID).Updates(transNueva).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error de inscripcion"), nil, c, http.StatusOK)
		return
	}
	_ = tx.Commit()
	utils.CrearRespuesta(nil, "Inscripcion exitosa", c, http.StatusOK)
}
