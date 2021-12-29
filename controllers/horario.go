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

func GetHorarios(c *gin.Context) {
	horarios := []*models.Horario{}
	var err error
	idParroquia := c.GetInt("id_parroquia")

	err = models.Db.Where(&models.Horario{ParroquiaID: uint(idParroquia)}).Preload("HorariosEntradas").Find(&horarios).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener horarios"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, horario := range horarios {
		if horario.Imagen == "" {
			horario.Imagen = utils.DefaultHorario
		} else {
			horario.Imagen = utils.SERVIMG + horario.Imagen
		}
		horario.Horarios = []string{}
		for _, entrada := range horario.HorariosEntradas {
			horario.Horarios = append(horario.Horarios, entrada.Descripcion)
		}

	}
	utils.CrearRespuesta(err, horarios, c, http.StatusOK)
}

func GetHorarioPorId(c *gin.Context) {
	horario := &models.Horario{}
	id := c.Param("id")
	err := models.Db.Preload("HorariosEntradas").First(horario, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Horario no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener horario"), nil, c, http.StatusInternalServerError)
		return
	}

	horario.Horarios = []string{}
	for _, entrada := range horario.HorariosEntradas {
		horario.Horarios = []string{}
		horario.Horarios = append(horario.Horarios, entrada.Descripcion)
	}
	if horario.Imagen == "" {
		horario.Imagen = utils.DefaultHorario
	} else {
		horario.Imagen = utils.SERVIMG + horario.Imagen
	}

	utils.CrearRespuesta(nil, horario, c, http.StatusOK)
}

func CreateHorario(c *gin.Context) {
	horario := &models.Horario{}
	err := c.ShouldBindJSON(horario)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Create(horario).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear horario"), nil, c, http.StatusInternalServerError)
		return
	}

	if horario.Imagen == "" {
		horario.Imagen = utils.DefaultHorario
	} else {
		idUrb := fmt.Sprintf("%d", horario.ID)
		horario.Imagen, err = img.FromBase64ToImage(horario.Imagen, "horarios/"+time.Now().Format(time.RFC3339)+idUrb, false)
		utils.Log.Info(horario.Imagen)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear horario "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Horario{}).Where("id = ?", horario.ID).Update("imagen", horario.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear horario "), nil, c, http.StatusInternalServerError)
			return
		}
		horario.Imagen = utils.SERVIMG + horario.Imagen
	}

	for _, hora := range horario.Horarios {
		err = tx.Create(&models.HorarioEntrada{HorarioID: horario.ID, Descripcion: hora}).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear horario "), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Horario creada con exito", c, http.StatusCreated)

}

func UpdateHorario(c *gin.Context) {
	horario := &models.Horario{}

	err := c.ShouldBindJSON(horario)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	idH, _ := strconv.Atoi(id)
	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(horario).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar horario"), nil, c, http.StatusInternalServerError)
		return
	}
	if img.IsBase64(horario.Imagen) {
		idUrb := fmt.Sprintf("%d", idH)
		horario.Imagen, err = img.FromBase64ToImage(horario.Imagen, "horarios/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear horario "), nil, c, http.StatusInternalServerError)
			return
		}
		err = tx.Model(&models.Horario{}).Where("id = ?", id).Update("imagen", horario.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar horario"), nil, c, http.StatusInternalServerError)
			return
		}
		horario.Imagen = utils.SERVIMG + horario.Imagen

	} else {
		horario.Imagen = utils.DefaultHorario
	}
	err = tx.Where("horario_id = ?", idH).Delete(&models.HorarioEntrada{}).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al crear horario "), nil, c, http.StatusInternalServerError)
		return
	}
	for _, hora := range horario.Horarios {
		err = tx.Create(&models.HorarioEntrada{HorarioID: uint(idH), Descripcion: hora}).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear horario "), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Horario actualizada correctamente", c, http.StatusOK)
}

func DeleteHorario(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Horario{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar horario"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Horario eliminada exitosamente", c, http.StatusOK)
}
