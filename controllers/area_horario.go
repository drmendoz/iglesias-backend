package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAreaHorarios(c *gin.Context) {
	horarios := []*models.AreaHorario{}
	err := models.Db.Find(&horarios).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener horarios"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, horario := range horarios {
		if horario.HoraInicio != nil && horario.HoraFin != nil {
			horario.HoraInicioSinFormato = fmt.Sprintf("%d:%d", horario.HoraInicio.Hour(), horario.HoraInicio.Minute())
			horario.HoraFinSinFormato = fmt.Sprintf("%d:%d", horario.HoraFin.Hour(), horario.HoraFin.Minute())
		}
	}
	utils.CrearRespuesta(err, horarios, c, http.StatusOK)
}

func GetAreaHorarioPorId(c *gin.Context) {
	horario := &models.AreaHorario{}
	id := c.Param("id")
	err := models.Db.First(horario, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("AreaHorario no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener horario"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, horario, c, http.StatusOK)
}

func CreateAreaHorario(c *gin.Context) {
	horario := &models.AreaHorario{}
	err := c.ShouldBindJSON(horario)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	if horario.HoraInicioSinFormato != "" {
		hora, err := strconv.Atoi(strings.Split(horario.HoraInicioSinFormato, ":")[0])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de apertura"), nil, c, http.StatusBadRequest)
			return
		}
		minutos, err := strconv.Atoi(strings.Split(horario.HoraInicioSinFormato, ":")[1])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de apertura"), nil, c, http.StatusBadRequest)
			return
		}
		tm := time.Date(1900, time.January, 0, hora, minutos, 0, 0, tiempo.Local)
		horario.HoraInicio = &tm
	}
	if horario.HoraFinSinFormato != "" {
		hora, err := strconv.Atoi(strings.Split(horario.HoraFinSinFormato, ":")[0])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de cierre"), nil, c, http.StatusBadRequest)
			return
		}
		minutos, err := strconv.Atoi(strings.Split(horario.HoraFinSinFormato, ":")[1])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de cierre"), nil, c, http.StatusBadRequest)
			return
		}
		tm := time.Date(1900, time.January, 0, hora, minutos, 0, 0, tiempo.Local)
		horario.HoraFin = &tm

	}

	err = models.Db.Create(horario).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear horario"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "Horario creado correctamente", c, http.StatusCreated)

}

func UpdateAreaHorario(c *gin.Context) {
	horario := &models.AreaHorario{}

	err := c.ShouldBindJSON(horario)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	if horario.HoraInicioSinFormato != "" {
		hora, err := strconv.Atoi(strings.Split(horario.HoraInicioSinFormato, ":")[0])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de apertura"), nil, c, http.StatusBadRequest)
			return
		}
		minutos, err := strconv.Atoi(strings.Split(horario.HoraInicioSinFormato, ":")[1])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de apertura"), nil, c, http.StatusBadRequest)
			return
		}
		tm := time.Date(1900, time.January, 0, hora, minutos, 0, 0, tiempo.Local)
		horario.HoraInicio = &tm
	}
	if horario.HoraFinSinFormato != "" {
		hora, err := strconv.Atoi(strings.Split(horario.HoraFinSinFormato, ":")[0])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de cierre"), nil, c, http.StatusBadRequest)
			return
		}
		minutos, err := strconv.Atoi(strings.Split(horario.HoraFinSinFormato, ":")[1])
		if err != nil {
			utils.CrearRespuesta(errors.New("Error en formato de hora de cierre"), nil, c, http.StatusBadRequest)
			return
		}
		tm := time.Date(1900, time.January, 0, hora, minutos, 0, 0, tiempo.Local)
		horario.HoraFin = &tm

	}

	err = models.Db.Where("id = ?", id).Updates(horario).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar horario"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(err, "AreaHorario actualizada correctamente", c, http.StatusOK)
}

func DeleteAreaHorario(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.AreaHorario{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar horario"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "AreaHorario eliminada exitosamente", c, http.StatusOK)
}

func GetHorariosPorAreaSocial(c *gin.Context) {
	idArea := c.Param("id")
	horarios := []*models.AreaHorario{}
	err := models.Db.Where("area_social_id = ?", idArea).Find(&horarios).Error
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al obtener horarios"), nil, c, 500)
		return

	}

	utils.CrearRespuesta(nil, horarios, c, 200)
}
