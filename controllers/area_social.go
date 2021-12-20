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

func GetAreaSocials(c *gin.Context) {

	etps := []*models.AreaSocial{}
	idParroquia := c.GetInt("id_etapa")
	err := models.Db.Where(&models.AreaSocial{ParroquiaID: uint(idParroquia)}).Order("created_at asc").Find(&etps).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener areas sociales"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, etp := range etps {
		if etp.Imagen == "" {
			etp.Imagen = utils.DefaultAreaSocial
		} else {
			etp.Imagen = utils.SERVIMG + etp.Imagen
		}
		if etp.ImagenReserva == "" {
			etp.ImagenReserva = utils.DefaultAreaSocial
		} else {
			etp.ImagenReserva = utils.SERVIMG + etp.ImagenReserva
		}
		etp.Estado = "Abierto"
		now := time.Now().In(tiempo.Local)
		fechaComparacion := time.Date(1900, time.January, 0, now.Hour(), now.Minute(), 0, 0, tiempo.Local)
		horarios := []*models.AreaHorario{}
		err = models.Db.Where("hora_inicio < ? ", fechaComparacion).Where("hora_fin > ?", fechaComparacion).Where("dia = ? ", now.Weekday().String()).Where("area_social_id = ?", etp.ID).First(&horarios).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				etp.Estado = "Cerrado"
			} else {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener area social"), nil, c, http.StatusInternalServerError)
				return
			}
		}
		etp.Horario = horarios
	}
	utils.CrearRespuesta(nil, etps, c, http.StatusOK)
}

func GetAreaSocialPorID(c *gin.Context) {
	etp := &models.AreaSocial{}
	id := c.Param("id")
	err := models.Db.Preload("Horarios", "fecha_fin >= ?", time.Now()).Preload("Reservaciones", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("hora_inicio ASC").Joins("Fiel").Joins("Fiel.Usuario")
	}).First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Area social no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener area"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultAreaSocial
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}
	if etp.ImagenReserva == "" {
		etp.ImagenReserva = utils.DefaultAreaSocial
	} else {
		etp.ImagenReserva = utils.SERVIMG + etp.ImagenReserva
	}
	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func GetAreaSocialPorId(c *gin.Context) {
	etp := &models.AreaSocial{}
	id := c.Param("id")
	fechaI := c.Query("fecha_inicio")
	fechaF := c.Query("fecha_fin")
	fechaInicio, _ := time.Parse("2006-01-02", fechaI)
	fechaFin, _ := time.Parse("2006-01-02", fechaF)
	err := models.Db.Preload("Horarios").Preload("Reservaciones", func(tx *gorm.DB) *gorm.DB {
		if fechaI == "" && fechaF == "" {
			return tx.Order("hora_inicio DESC").Joins("Fiel").Preload("Fiel.Usuario").Preload("Fiel.Casa")
		} else {
			return tx.Where("hora_inicio between ? and ? and valor_cancelado > ?", fechaInicio, fechaFin, 0).Order("hora_inicio DESC").Joins("Fiel").Preload("Fiel.Usuario").Preload("Fiel.Casa")
		}
	}).First(etp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Area social no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener area"), nil, c, http.StatusInternalServerError)
		return
	}
	if etp.Imagen == "" {
		etp.Imagen = utils.DefaultAreaSocial
	} else {
		etp.Imagen = utils.SERVIMG + etp.Imagen
	}
	if etp.ImagenReserva == "" {
		etp.ImagenReserva = utils.DefaultAreaSocial
	} else {
		etp.ImagenReserva = utils.SERVIMG + etp.ImagenReserva
	}
	etp.Estado = "Abierto"
	now := time.Now().In(tiempo.Local)
	fechaComparacion := time.Date(1900, time.January, 0, now.Hour(), now.Minute(), 0, 0, tiempo.Local)
	horarios := []*models.AreaHorario{}
	err = models.Db.Where("hora_inicio < ? ", fechaComparacion).Where("hora_fin > ?", fechaComparacion).Where("dia = ? ", now.Weekday().String(), etp.ID).Where("area_social_id = ?", etp.ID).Where("fecha_fin >= ?", time.Now()).First(&horarios).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			etp.Estado = "Cerrado"
		} else {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener area social"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	for i, hor := range etp.Horarios {
		etp.Horarios[i].HoraInicioSinFormato = fmt.Sprintf("%d:%d", hor.HoraInicio.Hour(), hor.HoraInicio.Minute())
		etp.Horarios[i].HoraFinSinFormato = fmt.Sprintf("%d:%d", hor.HoraFin.Hour(), hor.HoraFin.Minute())
	}
	etp.Horario = horarios
	utils.CrearRespuesta(nil, etp, c, http.StatusOK)
}

func CreateAreaSocial(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_etapa"))
	etp := &models.AreaSocial{}
	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	etp.ParroquiaID = idParroquia

	if etp.Imagen != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "areas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear area "), nil, c, http.StatusInternalServerError)

			return
		}
	}
	if etp.ImagenReserva != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.ImagenReserva, err = img.FromBase64ToImage(etp.ImagenReserva, "areas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear area "), nil, c, http.StatusInternalServerError)

			return
		}
	}

	tx := models.Db.Begin()
	err = tx.Create(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear area social"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, horario := range etp.Horario {
		horario.AreaSocialID = etp.ID
		err = tx.Create(horario).Error
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()

			utils.CrearRespuesta(errors.New("Error al obtener area social"), nil, c, http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Area social creada correctamente", c, http.StatusCreated)

}

func UpdateAreaSocial(c *gin.Context) {
	etp := &models.AreaSocial{}

	err := c.ShouldBindJSON(etp)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	if etp.Imagen != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.Imagen, err = img.FromBase64ToImage(etp.Imagen, "areas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al decodificar imagen"), nil, c, http.StatusInternalServerError)

			return
		}
	}
	if etp.ImagenReserva != "" {
		idUrb := fmt.Sprintf("%d", etp.ID)
		etp.ImagenReserva, err = img.FromBase64ToImage(etp.ImagenReserva, "areas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al decodificar imagen "), nil, c, http.StatusInternalServerError)

			return
		}
	}

	tx := models.Db.Begin()
	id := c.Param("id")
	err = tx.Where("id = ?", id).Updates(etp).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar area social"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Area social actualizada correctamente", c, http.StatusOK)
}

func DeleteAreaSocial(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.AreaSocial{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar area"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Area social eliminada exitosamente", c, http.StatusOK)
}

type HorarioDisponible struct {
	HoraInicio time.Time `json:"hora_inicio"`
	HoraFin    time.Time `json:"hora_fin"`
}

type AreaSocialDisponibles struct {
	AreaSocial         *models.AreaSocial   `json:"area_social"`
	HorarioDisponibles []*HorarioDisponible `json:"horarios_disponibles"`
	ExentoPago         bool                 `json:"exento_pago"`
}

func GetHorarioDisponiblesAreaSocial(c *gin.Context) {
	id := c.Param("id")
	fechaString := c.Query("fecha")
	fecha, err := time.Parse("2006-01-02", fechaString)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en formato de fecha"), nil, c, 400)
		return
	}
	dia := fecha.Weekday().String()
	area := &models.AreaSocial{}
	err = models.Db.First(area, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe area social"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener area social"), nil, c, http.StatusInternalServerError)
		return
	}
	fechaLimite := fecha.AddDate(0, 0, 1)
	horarios := []*HorarioDisponible{}
	if area.TiempoReservacionMinutos == 0 {
		utils.CrearRespuesta(errors.New("Error al obtener horarios comunicarse con el administrador"), nil, c, 500)
		return
	}
	for fecha.Before(fechaLimite) {

		fechaHoraInicio := time.Date(1900, time.January, 0, fecha.Hour(), fecha.Minute(), 0, 0, tiempo.Local)
		fechaHoraFin := fechaHoraInicio.Add(time.Duration(area.TiempoReservacionMinutos * int(time.Minute)))
		var result int64
		err = models.Db.Model(&models.AreaHorario{}).Where("hora_inicio <= ? ", fechaHoraInicio).Where("hora_fin >= ?", fechaHoraFin).Where("dia = ? ", dia).Where("area_social_id = ?", id).Where("fecha_fin >= ?", time.Now()).Count(&result).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error interno del servidor"), nil, c, http.StatusInternalServerError)
			return
		}
		if result > 0 {
			var resultReservacion int64
			fechaInicio := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), fechaHoraInicio.Hour(), fechaHoraInicio.Minute(), 0, 0, tiempo.Local)
			fechaFin := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), fechaHoraFin.Hour(), fechaHoraFin.Minute(), 0, 0, tiempo.Local)
			err = models.Db.Model(&models.ReservacionAreaSocial{}).Where("hora_inicio <= ? ", fechaInicio).Where("hora_fin >= ?", fechaFin).Where("area_social_id = ?", id).Count(&resultReservacion).Error
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error interno del servidor"), nil, c, http.StatusInternalServerError)
				return
			}
			print("Conteo reservacion :", resultReservacion)
			if resultReservacion == 0 {
				horarios = append(horarios, &HorarioDisponible{HoraInicio: fechaInicio, HoraFin: fechaFin})
			}
		}
		fecha = fecha.Add(time.Duration(area.TiempoReservacionMinutos * int(time.Minute)))
	}
	if area.Imagen == "" {
		area.Imagen = utils.DefaultAreaSocial
	} else {
		area.Imagen = utils.SERVIMG + area.Imagen
	}
	exentoPago := c.GetInt("id_residente") == 34
	areaDisponbiles := AreaSocialDisponibles{AreaSocial: area, HorarioDisponibles: horarios, ExentoPago: exentoPago}
	utils.CrearRespuesta(nil, areaDisponbiles, c, http.StatusOK)
}
