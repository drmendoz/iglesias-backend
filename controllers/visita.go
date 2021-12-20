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
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/notification"
	"github.com/drmendoz/iglesias-backend/utils/tiempo"
	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

func GetVisitas(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	idCasa := uint(c.GetInt("id_casa"))
	tipoUsuario := c.Query("tipo_usuario")
	mz := c.Query("mz")
	villa := c.Query("villa")
	fecha := c.Query("fecha")
	idAutorizacion, err := strconv.Atoi(c.Query("id_autorizacion"))
	if err != nil && idAutorizacion != 0 {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener visitas"), nil, c, http.StatusInternalServerError)
		return
	}
	fechaCreacion, err := time.Parse(time.RFC3339, fecha)
	if err != nil && idAutorizacion != 0 {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener visitas"), nil, c, http.StatusInternalServerError)
		return
	}
	visitas := []*models.Visita{}
	if idCasa == 0 && (mz != "" && villa != "") {
		casa := &models.Casa{}
		err = models.Db.Model(&models.Casa{}).Where(&models.Casa{Manzana: mz, Villa: villa, EtapaID: idEtapa}).First(casa).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener casa"), nil, c, http.StatusInternalServerError)
			return
		}
		idCasa = casa.ID
	}
	idAut := uint(idAutorizacion)
	busqueda := &models.Visita{}
	if idAut != 0 {
		busqueda = &models.Visita{CasaID: idCasa, EtapaID: idEtapa, TipoUsuario: tipoUsuario, AutorizacionID: &idAut, DiaCreacion: fechaCreacion}
	} else {
		busqueda = &models.Visita{CasaID: idCasa, EtapaID: idEtapa, TipoUsuario: tipoUsuario, DiaCreacion: fechaCreacion}
	}
	if idCasa == 0 && (mz != "" && villa == "") {
		err = models.Db.Limit(100).Where(busqueda).Where("Casa.manzana = ?", mz).Order("updated_at desc").Joins("Publicador").Joins("Usuario").Joins("Casa").Joins("Autorizacion").Joins("Entrada").Joins("Salida").Find(&visitas).Error
	} else {
		err = models.Db.Limit(100).Where(busqueda).Order("updated_at desc").Joins("Publicador").Joins("Usuario").Joins("Casa").Joins("Autorizacion").Joins("Entrada").Joins("Salida").Find(&visitas).Error
	}

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener visitas"), nil, c, http.StatusInternalServerError)
		return
	}
	tx := models.Db.Begin()
	for _, visita := range visitas {
		tiempoMax := 60
		visita.SegundosRestantes = time.Duration(tiempoMax) - time.Duration(time.Since(visita.CreatedAt).Seconds())
		visita.SegundosTotal = tiempoMax
		if uint(visita.SegundosRestantes) > uint(tiempoMax) && visita.Estado == "PENDIENTE" {
			visita.Estado = "ESPERANDO"
			err = tx.Omit("imagen").Where("id = ?", visita.ID).Updates(visita).Error
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
				return
			}
		}
		if visita.Nuevo {
			err = tx.Model(&models.Visita{}).Omit("updated_at").Where("id = ?", visita.ID).Update("nuevo", false).Error
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
				return
			}
		}

		visita.Publicador.Contrasena = ""
		if visita.UsuarioID != 0 {

			visita.Usuario.Contrasena = ""
		}
		if visita.Imagen == "" {
			visita.Imagen = utils.DefaultVisita
		} else {
			visita.Imagen = utils.SERVIMG + visita.Imagen
		}
		if visita.Imagen2 == "" {
			visita.Imagen2 = utils.DefaultVisita
		} else {
			visita.Imagen2 = utils.SERVIMG + visita.Imagen
		}
		if visita.RespuestaPorLlamada {
			visita.Usuario = &models.Usuario{Nombre: "llamada"}
		}
		if visita.Entrada != nil {
			visita.Entrada.TipoEntrada = strcase.ToCamel(strings.ToLower(visita.Entrada.TipoEntrada))
		}
		visita.TipoUsuario = strcase.ToCamel(strings.ToLower(visita.TipoUsuario))
		visita.TipoEntrada = strcase.ToCamel(strings.ToLower(visita.TipoEntrada))

	}
	tx.Commit()

	utils.CrearRespuesta(err, visitas, c, http.StatusOK)
}

func GetVisitasPorEtapa(c *gin.Context) {
	idEtapa := c.Param("id")

	idE, err := strconv.Atoi(idEtapa)
	if err != nil {
		utils.CrearRespuesta(errors.New("Formato de id etapa incorrecto"), nil, c, http.StatusBadRequest)
		return
	}
	visitas := []*models.Visita{}
	err = models.Db.Where(&models.Visita{EtapaID: uint(idE)}).Order("created_at desc").Joins("Etapa").Joins("Publicador").Joins("Usuario").Joins("Casa").Joins("Entrada").Joins("Salida").Find(&visitas).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener visitas"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, visita := range visitas {
		visita.Publicador.Contrasena = ""
		if visita.UsuarioID != 0 {

			visita.Usuario.Contrasena = ""
		}
		if visita.Imagen == "" {
			visita.Imagen = utils.DefaultVisita
		} else {
			visita.Imagen = utils.SERVIMG + visita.Imagen
		}
		if visita.Imagen2 == "" {
			visita.Imagen2 = utils.DefaultVisita
		} else {
			visita.Imagen2 = utils.SERVIMG + visita.Imagen
		}

	}
	utils.CrearRespuesta(err, visitas, c, http.StatusOK)
}

func GetVisitaPorId(c *gin.Context) {
	visita := &models.Visita{}
	id := c.Param("id")
	err := models.Db.Joins("Publicador").Joins("Usuario").Joins("Entrada").Joins("Salida").First(visita, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Visita no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener visita"), nil, c, http.StatusInternalServerError)
		return
	}
	if visita.RespuestaPorLlamada {
		visita.Usuario = &models.Usuario{Nombre: "llamada"}
	}
	if visita.Imagen == "" {
		visita.Imagen = utils.DefaultVisita
	} else {
		visita.Imagen = utils.SERVIMG + visita.Imagen
	}
	if visita.Imagen2 == "" {
		visita.Imagen2 = utils.DefaultVisita
	} else {
		visita.Imagen2 = utils.SERVIMG + visita.Imagen
	}
	visita.Publicador.Contrasena = ""
	if visita.UsuarioID != 0 {
		visita.Usuario.Contrasena = ""
	}
	tiempoMax := 70
	visita.SegundosRestantes = time.Duration(tiempoMax) - time.Duration(time.Since(visita.CreatedAt).Seconds())
	visita.SegundosTotal = tiempoMax
	if uint(visita.SegundosRestantes) > uint(tiempoMax) && visita.Estado == "PENDIENTE" {
		visita.Estado = "ESPERANDO"
		err := models.Db.Where("id = ?", id).Updates(visita).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	utils.CrearRespuesta(nil, visita, c, http.StatusOK)
}

func CreateVisita(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	idEtapa := c.GetInt("id_etapa")
	if idUsuario == 0 || idEtapa == 0 {
		utils.CrearRespuesta(errors.New("Error al crear visita"), nil, c, http.StatusInternalServerError)
		return
	}
	visita := &models.Visita{}

	err := c.ShouldBindJSON(visita)
	visita.PublicadorID = uint(idUsuario)
	visita.EtapaID = uint(idEtapa)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}

	if visita.Entrada != nil {
		visita.Entrada.HoraEntrada = time.Now()
	} else if visita.Entrada == nil && visita.Salida != nil {
		visita.SalidaFirst = true
	}
	t := time.Now()
	visita.DiaCreacion = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	idUrb := fmt.Sprintf("%d", visita.ID)
	if visita.Imagen != "" {
		visita.Imagen, err = img.FromBase64ToImage(visita.Imagen, "visitas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear visita "), nil, c, http.StatusInternalServerError)

			return
		}
	}
	if visita.Imagen2 != "" {
		visita.Imagen2, err = img.FromBase64ToImage(visita.Imagen2, "visitas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear visita "), nil, c, http.StatusInternalServerError)

			return
		}
	}
	tx := models.Db.Begin()
	err = tx.Omit("usuario_id").Create(visita).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al crear visita "), nil, c, http.StatusInternalServerError)

		return
	}
	residentes := []*models.Residente{}
	err = tx.Where("casa_id = ?", visita.CasaID).Find(&residentes).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(err, "Error al crear publicacion", c, http.StatusInternalServerError)
		return
	}
	tokens := []string{}
	for _, res := range residentes {
		tokens = append(tokens, res.TokenNotificacion)
	}
	go notification.SendNotification("Nuevo visitante: "+visita.Nombre, "Motivo: "+visita.Motivo, tokens, "4")
	tx.Commit()
	utils.CrearRespuesta(err, visita.ID, c, http.StatusCreated)
}

func CreateSalidaVisita(visita *models.Visita, tx *gorm.DB) (*models.SalidaVisita, error) {
	salida := &models.SalidaVisita{}
	salida.HoraSalida = time.Now()
	salida.TipoSalida = visita.Salida.TipoSalida
	salida.Placa = visita.Salida.Placa
	err := models.Db.Create(salida).Error
	if err != nil {
		return salida, err
	} else {
		return salida, nil
	}
}

func UpdateVisita(c *gin.Context) {
	visita := &models.Visita{}

	err := c.ShouldBindJSON(visita)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")

	tx := models.Db.Begin()
	err = tx.Omit("imagen").Where("id = ?", id).Updates(visita).Error
	if err != nil {
		_ = c.Error(err)
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
		return
	}
	if visita.Imagen != "" {
		idUrb := fmt.Sprintf("%d", visita.ID)
		visita.Imagen, err = img.FromBase64ToImage(visita.Imagen, "visitas/"+time.Now().Format(time.RFC3339)+idUrb, false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear visita "), nil, c, http.StatusInternalServerError)

			return
		}
		err = tx.Model(&models.Visita{}).Where("id = ?", visita.ID).Update("imagen", visita.Imagen).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
			return
		}
		visita.Imagen = utils.SERVIMG + visita.Imagen

	} else {
		visita.Imagen = utils.DefaultVisita
	}

	if visita.Salida != nil {
		salida, err := CreateSalidaVisita(visita, tx)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al ingresar salida"), nil, c, http.StatusInternalServerError)
			return
		}
		visita.SalidaID = &salida.ID
		err = tx.Model(&models.Visita{}).Where("id = ?", visita.ID).Update("salida_id", salida.ID).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Visita actualizada correctamente", c, http.StatusOK)
}

func DeleteVisita(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Visita{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar visita"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Visita eliminada exitosamente", c, http.StatusOK)
}

type NotificacionGarita struct {
	Show              bool          `json:"show"`
	Mensaje           string        `json:"mensaje"`
	TipoNotificacion  string        `json:"tipo_notificacion"`
	Celular           string        `json:"celular"`
	SegundosRestantes time.Duration `json:"seconds_remaining" gorm:"-"`
	SegundosTotal     int           `json:"seconds_total" gorm:"-"`
}

func NotificarVisita(c *gin.Context) {
	idUsuarioGarita := c.GetInt("id_usuario")
	notificacion := &NotificacionGarita{Show: false}
	visita := &models.Visita{}
	err := models.Db.Where("publicador_id = ?", idUsuarioGarita).Joins("Casa").Last(visita).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe entrada"), nil, c, http.StatusNotFound)
			return
		}
		utils.CrearRespuesta(errors.New("Error interno del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	if visita.Estado == "ACEPTADA" && !visita.Vista {

		notificacion.Show = true
		notificacion.Mensaje = "Visita aceptada"
		notificacion.TipoNotificacion = "exito"

		_ = models.Db.Updates(&models.Visita{Vista: true})
	} else if visita.Estado == "RECHAZADA" && !visita.Vista {

		notificacion.Show = true
		notificacion.Mensaje = "Visita rechazada"

		notificacion.TipoNotificacion = "falla"

		_ = models.Db.Updates(&models.Visita{Vista: true})
	} else if visita.CreatedAt.Add(time.Second*30).After(time.Now().In(tiempo.Local)) && !visita.Vista {
		notificacion.Show = true

		_ = models.Db.Updates(&models.Visita{Vista: true})
		notificacion.Mensaje = "No hay respuesta de la casa. Por favor llamar a  residente."
		notificacion.TipoNotificacion = "pendiente"

		notificacion.Celular = visita.Casa.Celular
	}
	utils.CrearRespuesta(nil, notificacion, c, http.StatusOK)
}

func NotificarVisitaPorId(c *gin.Context) {
	id := c.Param("id")
	notificacion := &NotificacionGarita{Show: false}
	visita := &models.Visita{}
	err := models.Db.Joins("Casa").First(visita, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe entrada"), nil, c, http.StatusNotFound)
			return
		}
		utils.CrearRespuesta(errors.New("Error interno del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	tiempoMax := 10
	visita.SegundosRestantes = time.Duration(tiempoMax) - time.Duration(time.Since(visita.CreatedAt).Seconds())
	visita.SegundosTotal = tiempoMax
	if uint(visita.SegundosRestantes) > uint(tiempoMax) && visita.Estado == "PENDIENTE" {
		visita.Estado = "ESPERANDO"
		err = models.Db.Omit("imagen").Where("id = ?", visita.ID).Updates(visita).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar visita"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	notificacion.SegundosRestantes = visita.SegundosRestantes
	notificacion.SegundosTotal = visita.SegundosTotal

	if visita.Estado == "ACEPTADA" && !visita.Vista {

		notificacion.Show = true
		notificacion.Mensaje = "Visita aceptada"
		notificacion.TipoNotificacion = "exito"

		_ = models.Db.Updates(&models.Visita{Vista: true})
	} else if visita.Estado == "RECHAZADA" && !visita.Vista {

		notificacion.Show = true
		notificacion.Mensaje = "Visita rechazada"

		notificacion.TipoNotificacion = "falla"

		_ = models.Db.Updates(&models.Visita{Vista: true})
	} else if visita.CreatedAt.Add(time.Second*30).After(time.Now().In(tiempo.Local)) && !visita.Vista {
		notificacion.Show = true

		_ = models.Db.Updates(&models.Visita{Vista: true})
		notificacion.Mensaje = "No hay respuesta de la casa. Por favor llamar a  residente."
		notificacion.TipoNotificacion = "pendiente"

		notificacion.Celular = visita.Casa.Celular
	}

	utils.CrearRespuesta(nil, notificacion, c, http.StatusOK)
}
func ContestarVisita(c *gin.Context) {
	idVisita := c.Param("id")
	idUsuario := c.GetInt("id_usuario")
	visita := &models.Visita{}
	err := c.ShouldBindJSON(visita)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al validar parametros"), nil, c, http.StatusNotAcceptable)
		return
	}
	if idUsuario == 0 {
		utils.CrearRespuesta(errors.New("Usuario no identificado"), nil, c, http.StatusNotAcceptable)
		return
	}
	visita.UsuarioID = uint(idUsuario)
	err = models.Db.Where("id = ?", idVisita).Select("estado", "usuario_id").Updates(visita).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al responder visita"), nil, c, http.StatusInternalServerError)
		return
	}

	utils.CrearRespuesta(nil, "Solicitud enviada con exito", c, http.StatusOK)
}
