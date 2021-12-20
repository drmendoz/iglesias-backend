package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/notification"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetVotacions(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	idFiel := c.GetInt("id_residente")
	votacions := []*models.Votacion{}
	var err error
	if idParroquia != 0 {
		err = models.Db.Order("created_at DESC").Where("etapa_id = ?", idParroquia).Preload("Opciones", func(db *gorm.DB) *gorm.DB {
			return db.Order("opcion_votacion.conteo DESC")
		}).Find(&votacions).Error
	} else {

		err = models.Db.Preload("Opciones").Find(&votacions).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener votacions"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, vot := range votacions {
		voto := false
		if idFiel != 0 {
			res := &models.RespuestaVotacion{}
			err = models.Db.Where("residente_id = ? and OpcionVotacion.votacion_id = ? ", idFiel, vot.ID).Joins("OpcionVotacion").First(res).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error al obtener votaciones"), nil, c, http.StatusInternalServerError)
					return
				}
			} else {
				voto = true
			}
			utils.Log.Info(c.GetBool("is_principal"))
			if !c.GetBool("is_principal") {
				voto = true
			}
		}
		vot.Expiro = vot.FechaVencimiento.Before(time.Now().In(tiempo.Local))
		if vot.Expiro || voto {
			vot.UsuarioVotacion = true
		}
		for i, opc := range vot.Opciones {
			opc.Color = utils.Colores[i+1]
			if int(vot.TotalVotos) != 0 {

				opc.Porcentaje = float64(opc.Conteo) / float64(vot.TotalVotos)
			}
		}
		imagenesArr := strings.Split(vot.Imagenes, ",")
		imagenes := []string{}
		if vot.Imagenes != "" {
			for _, imagen := range imagenesArr {
				imagen = utils.SERVIMG + imagen
				imagenes = append(imagenes, imagen)
			}
		} else {
			imagenes = append(imagenes, utils.DefaultVotacion)
		}
		vot.ImagenesArray = imagenes
	}
	utils.CrearRespuesta(nil, votacions, c, http.StatusOK)
}

func GetVotacionPorId(c *gin.Context) {
	usuarioVoto := c.GetBool("usuario_voto")
	votacion := &models.Votacion{}
	id := c.Param("id")
	err := models.Db.Preload("Opciones", func(db *gorm.DB) *gorm.DB {
		return db.Order("opcion_votacion.conteo DESC")
	}).First(votacion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Votacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener votacion"), nil, c, http.StatusInternalServerError)
		return
	}
	votacion.Expiro = votacion.FechaVencimiento.Before(time.Now().In(tiempo.Local))
	if usuarioVoto || votacion.Expiro {
		votacion.UsuarioVotacion = true
	}
	for i, opc := range votacion.Opciones {
		opc.Color = utils.Colores[i+1]
		if int(votacion.TotalVotos) != 0 {

			opc.Porcentaje = float64(opc.Conteo) / float64(votacion.TotalVotos)
		}
	}
	imagenesArr := strings.Split(votacion.Imagenes, ",")
	imagenes := []string{}
	if votacion.Imagenes != "" {
		for _, imagen := range imagenesArr {
			imagen = utils.SERVIMG + imagen
			imagenes = append(imagenes, imagen)
		}
	} else {
		imagenes = append(imagenes, utils.DefaultVotacion)
	}
	votacion.ImagenesArray = imagenes
	utils.CrearRespuesta(nil, votacion, c, http.StatusOK)
}

func CreateVotacion(c *gin.Context) {
	idParroquia := c.GetInt("id_etapa")
	if idParroquia == 0 {
		utils.CrearRespuesta(errors.New("No existe el id_etapa"), nil, c, http.StatusOK)
		return
	}

	votacion := &models.Votacion{}
	err := c.ShouldBindJSON(votacion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	votacion.ParroquiaID = uint(idParroquia)

	if !votacion.FechaVencimiento.IsZero() {
		rounded := time.Date(votacion.FechaVencimiento.Year(), votacion.FechaVencimiento.Month(), votacion.FechaVencimiento.Day(), 0, 0, 0, 0, time.Local)
		votacion.FechaVencimiento = rounded.Add(24 * time.Hour)
	}

	imagenesArr := votacion.ImagenesArray
	if len(imagenesArr) > 0 {
		idUrb := fmt.Sprintf("%d", votacion.ParroquiaID)
		imagenes := []string{}
		for _, imagen := range imagenesArr {
			imagen, err = img.FromBase64ToImage(imagen, "votaciones/"+time.Now().Format(time.RFC3339Nano)+idUrb, false)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear votacion "), nil, c, http.StatusInternalServerError)
				return
			}
			imagenes = append(imagenes, imagen)
		}
		votacion.Imagenes = strings.Join(imagenes, ",")
	}

	err = tx.Omit("Opciones").Create(votacion).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear votacion"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, opcion := range votacion.Opciones {
		opcion.VotacionID = votacion.ID
		err = tx.Create(&opcion).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear votacion"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	residentes := []*models.Fiel{}
	err = models.Db.Where("Casa.etapa_id = ?", votacion.ParroquiaID).Joins("Casa").Find(&residentes).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(err, "Error al crear publicacion", c, http.StatusInternalServerError)
		return
	}
	tokens := []string{}
	for _, res := range residentes {
		tokens = append(tokens, res.TokenNotificacion)
	}

	go notification.SendNotification("Nueva Encuesta disponible", "Pregunta: "+votacion.Pregunta, tokens, "2")
	utils.CrearRespuesta(err, "Votacion creada correctamente", c, http.StatusCreated)

}

func UpdateVotacion(c *gin.Context) {
	votacion := &models.Votacion{}

	err := c.ShouldBindJSON(votacion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	err = tx.Omit("Opciones", "imagen").Where("id = ?", id).Updates(votacion).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar votacion"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, opcion := range votacion.Opciones {
		opcion.VotacionID = votacion.ID
		err = tx.Updates(&opcion).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar votacion"), nil, c, http.StatusInternalServerError)
		}
	}
	imagenesArr := votacion.ImagenesArray
	if len(imagenesArr) > 0 {
		idUrb := fmt.Sprintf("%d", votacion.ParroquiaID)
		imagenes := []string{}
		for _, imagen := range imagenesArr {
			imagen, err = img.FromBase64ToImage(imagen, "votaciones/"+time.Now().Format(time.RFC3339Nano)+idUrb, false)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear votacion "), nil, c, http.StatusInternalServerError)
				return
			}
			imagenes = append(imagenes, imagen)
		}
		err = tx.Model(&models.Votacion{}).Where("id = ?", id).Update("imagenes_string", strings.Join(imagenes, ",")).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear votacion "), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Votacion actualizada correctamente", c, http.StatusOK)
}

func DeleteVotacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Votacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar votacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Votacion eliminada exitosamente", c, http.StatusOK)
}

func ResponderVotacion(c *gin.Context) {
	usuarioVoto := c.GetBool("usuario_voto")
	idUsuario := c.GetInt("id_residente")
	idVotacion := c.Param("id")
	res := &models.RespuestaVotacion{}
	err := c.ShouldBindJSON(&res)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en los parametros de la respuesta"), nil, c, http.StatusBadRequest)
		return
	}
	if usuarioVoto {
		utils.CrearRespuesta(errors.New("Usuario ya voto en esta encuesta"), nil, c, http.StatusNotAcceptable)
		return
	}
	res.FielID = uint(idUsuario)
	tx := models.Db.Begin()
	err = tx.Create(&res).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al ingresar voto"), nil, c, http.StatusOK)
		return
	}
	err = tx.Model(&models.OpcionVotacion{}).Where("id = ?", res.OpcionVotacionID).Update("conteo", gorm.Expr("conteo + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al ingresar voto"), nil, c, http.StatusOK)
		return
	}

	err = tx.Model(&models.Votacion{}).Where("id = ?", idVotacion).Update("total_votos", gorm.Expr("total_votos + ?", 1)).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al ingresar voto"), nil, c, http.StatusOK)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Votacion exitosa", c, http.StatusOK)
}
