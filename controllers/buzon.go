package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/slice"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEntradasBuzon(c *gin.Context) {
	buzones := []*models.Buzon{}

	err := models.Db.Preload("Publicador").Preload("Destinatarios").Find(&buzones).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, buzon := range buzones {
		buzon.UltimoMensaje = &models.Buzon{}
		err = models.Db.First(buzon.UltimoMensaje, "buzon_remitente_id = ?", buzon.BuzonRemitenteID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				buzon.UltimoMensaje = nil
				continue
			}
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
			return
		}
	}
	slice.Sort(buzones[:], func(i, j int) bool {
		return buzones[i].UltimoMensaje != nil && buzones[i].UltimoMensaje.CreatedAt.Before(buzones[j].UltimoMensaje.CreatedAt)
	})
	utils.CrearRespuesta(nil, buzones, c, http.StatusOK)
}

func GetBuzonesEnviados(c *gin.Context) {
	buzones := []*models.Buzon{}
	idUsuario := c.GetInt("id_usuario")
	err := models.Db.Where("publicador_id = ?", idUsuario).Preload("Publicador").Preload("Destinatarios").Preload("Destinatarios.Casa").Order("created_at desc").Preload("Casa").Find(&buzones).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, buzon := range buzones {
		buzon.Leido = true
		buzon.UltimoMensaje = &models.Buzon{}
		err = models.Db.Last(buzon.UltimoMensaje, "buzon_remitente_id = ?", buzon.BuzonRemitenteID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				buzon.UltimoMensaje = nil
				continue
			}
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
			return
		}

	}
	slice.Sort(buzones[:], func(i, j int) bool {
		if buzones[i].UltimoMensaje == nil {
			buzones[i].UltimoMensaje = &models.Buzon{}
			buzones[i].UltimoMensaje.CreatedAt = buzones[i].CreatedAt
			buzones[i].UltimoMensaje.Titulo = buzones[i].Titulo
			buzones[i].UltimoMensaje.Descripcion = buzones[i].Descripcion
		}
		if buzones[j].UltimoMensaje == nil {
			buzones[j].UltimoMensaje = &models.Buzon{}
			buzones[j].UltimoMensaje.CreatedAt = buzones[j].CreatedAt
		}
		return buzones[i].UltimoMensaje.CreatedAt.After(buzones[j].UltimoMensaje.CreatedAt)
	})
	slice.Sort(buzones[:], func(i, j int) bool {
		return !buzones[i].Leido
	})
	utils.CrearRespuesta(nil, buzones, c, http.StatusOK)
}

func GetBuzonesRecibidosAdminParroquia(c *gin.Context) {
	buzones := []*models.Buzon{}
	idUsuario := c.GetInt("id_usuario")
	err := models.Db.Where("is_admin = ?", false).Where("buzon_remitente_id is not null ").Preload("Publicador").Preload("Destinatarios.Casa").Preload("Destinatarios").Preload("Casa").Order("created_at desc").Find(&buzones).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
		return
	}
	buzonesRecibidos := []*models.Buzon{}
	for _, buzon := range buzones {
		bandera := false
		for _, b := range buzonesRecibidos {
			buzonRPub := buzon.PublicadorID
			bRPub := b.PublicadorID
			buzonRRem := *buzon.BuzonRemitenteID
			bRRem := *b.BuzonRemitenteID
			if bRPub == buzonRPub && buzonRRem == bRRem {
				bandera = true
			}
		}
		if bandera {
			continue
		}
		buzon.UltimoMensaje = &models.Buzon{}
		err = models.Db.Last(buzon.UltimoMensaje, "buzon_remitente_id = ?", buzon.BuzonRemitenteID).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
				return
			} else {
				buzon.UltimoMensaje = nil
			}
		}
		var result int64
		err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: buzon.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
			return
		}
		buzon.Leido = result > 0
		if !buzon.Leido {
			var result int64
			err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: buzon.UltimoMensaje.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
				return
			}
			buzon.Leido = result > 0
		}
		if buzon.BuzonRemitenteID != nil {
			destinatarios := []*models.BuzonDestinatario{}
			err = models.Db.Where("buzon_id = ?", buzon.ID).Preload("Casa").Find(&destinatarios).Error
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
				return
			}
			buzon.Destinatarios = destinatarios
		}
		buzonesRecibidos = append(buzonesRecibidos, buzon)
	}
	slice.Sort(buzonesRecibidos[:], func(i, j int) bool {
		if buzonesRecibidos[i].UltimoMensaje == nil {
			buzonesRecibidos[i].UltimoMensaje = &models.Buzon{}
			buzonesRecibidos[i].UltimoMensaje.CreatedAt = buzonesRecibidos[i].CreatedAt
			buzonesRecibidos[i].UltimoMensaje.Titulo = buzonesRecibidos[i].Titulo
			buzonesRecibidos[i].UltimoMensaje.Descripcion = buzonesRecibidos[i].Descripcion
		}
		return buzonesRecibidos[i].UltimoMensaje.CreatedAt.After(buzonesRecibidos[j].UltimoMensaje.CreatedAt)
	})
	slice.Sort(buzonesRecibidos[:], func(i, j int) bool {
		return !buzonesRecibidos[i].Leido
	})

	utils.CrearRespuesta(nil, buzonesRecibidos, c, http.StatusOK)
}

type Buzones struct {
	Enviados  []*models.Buzon `json:"enviados"`
	Recibidos []*models.Buzon `json:"recibidos"`
}

func GetBuzonesResidente(c *gin.Context) {
	id := c.GetInt("id_usuario")
	idCasa := c.GetInt("id_casa")
	idUsuario := uint(id)
	buzones := &Buzones{}
	buzonesRecibidos, err := getBuzonesRecibidosResidente(int(idUsuario), idCasa)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzon"), nil, c, http.StatusInternalServerError)
		return
	}
	buzonesEnviados, err := getBuzonesEnviadosResidente(int(idUsuario))
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzon"), nil, c, http.StatusInternalServerError)
		return
	}
	buzones.Recibidos = buzonesRecibidos
	buzones.Enviados = buzonesEnviados

	utils.CrearRespuesta(nil, buzones, c, http.StatusOK)
}

func getBuzonesRecibidosResidente(idUsuario int, idCasa int) ([]*models.Buzon, error) {
	recibidos := []*models.Buzon{}
	publicos := []*models.Buzon{}
	err := models.Db.Where("publico = ?", true).Where("publicador_id != ?", idUsuario).Where("is_admin = ?", true).Preload("Publicador").Preload("Destinatarios").Preload("Archivos").Preload("Mensajes").Preload("Mensajes.Archivos").Order("created_at desc").Find(&publicos).Error
	if err != nil {
		return nil, err
	}
	for _, buzon := range publicos {

		buzon.UltimoMensaje = &models.Buzon{}
		err = models.Db.First(buzon.UltimoMensaje, "buzon_remitente_id = ?", buzon.BuzonRemitenteID).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			} else {
				buzon.UltimoMensaje = nil
			}
		}

		var result int64
		err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: buzon.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
		if err != nil {
			return nil, err
		}
		buzon.Leido = result > 0
		for _, arc := range buzon.Archivos {
			arc.Url = utils.SERVIMG + arc.Url
		}
		buzon.Adjuntos = len(buzon.Archivos) > 0
		for _, msj := range buzon.Mensajes {
			var result int64
			err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: msj.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
			if err != nil {
				return nil, err
			}
			msj.Leido = result > 0
			if msj.Leido {
				buzon.Leido = true
			}
			if len(msj.Archivos) > 0 {
				buzon.Adjuntos = true
			}
		}
		recibidos = append(recibidos, buzon)
	}
	destinatarios := []*models.BuzonDestinatario{}
	err = models.Db.Where("casa_id = ?", idCasa).Preload("Buzon").Preload("Buzon.Publicador").Preload("Buzon.Destinatarios").Preload("Buzon.Archivos").Preload("Buzon.Mensajes").Preload("Buzon.Mensajes.Archivos").Order("created_at desc").Find(&destinatarios).Error
	if err != nil {
		return nil, err
	}
	for _, dest := range destinatarios {
		if dest.Buzon != nil {
			dest.Buzon.UltimoMensaje = &models.Buzon{}
			err = models.Db.First(dest.Buzon.UltimoMensaje, "buzon_remitente_id = ?", dest.Buzon.BuzonRemitenteID).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				} else {
					dest.Buzon.UltimoMensaje = nil
				}

			}
			mensajes := []*models.Buzon{}
			for _, msj := range dest.Buzon.Mensajes {
				if msj.IsAdmin || msj.PublicadorID == uint(idUsuario) {
					mensajes = append(mensajes, msj)
				}
			}
			dest.Buzon.Mensajes = mensajes
			var result int64
			err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: dest.Buzon.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
			if err != nil {
				return nil, err
			}
			dest.Buzon.Leido = result > 0
			for _, arc := range dest.Buzon.Archivos {
				arc.Url = utils.SERVIMG + arc.Url
			}
			dest.Buzon.Adjuntos = len(dest.Buzon.Archivos) > 0
			for _, msj := range dest.Buzon.Mensajes {
				var result int64
				err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: msj.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
				if err != nil {
					return nil, err
				}
				msj.Leido = result > 0
				if msj.Leido {
					dest.Buzon.Leido = true
				}
				if len(msj.Archivos) > 0 {
					dest.Buzon.Adjuntos = true
				}
			}
			recibidos = append(recibidos, dest.Buzon)

		}
	}
	slice.Sort(recibidos[:], func(i, j int) bool {
		return recibidos[i].CreatedAt.After(recibidos[j].CreatedAt)
	})
	return recibidos, nil

}

func getBuzonesEnviadosResidente(idUsuario int) ([]*models.Buzon, error) {
	buzonesEnviados := []*models.Buzon{}
	err := models.Db.Where("publicador_id = ?", idUsuario).Preload("Publicador").Preload("Destinatarios").Preload("Archivos").Order("created_at desc").Find(&buzonesEnviados).Error
	if err != nil {
		return nil, err
	}
	for _, buzon := range buzonesEnviados {
		buzon.Leido = true
		buzon.UltimoMensaje = &models.Buzon{}
		err = models.Db.First(buzon.UltimoMensaje, "buzon_remitente_id = ?", buzon.BuzonRemitenteID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				buzon.UltimoMensaje = nil
				continue
			}
		}
		for _, arc := range buzon.Archivos {
			arc.Url = utils.SERVIMG + arc.Url
		}
	}
	slice.Sort(buzonesEnviados[:], func(i, j int) bool {
		return buzonesEnviados[i].CreatedAt.After(buzonesEnviados[j].CreatedAt)
	})
	return buzonesEnviados, nil
}
func GetMensajesBuzones(c *gin.Context) {
	rol := c.GetString("rol")
	idUsuario := c.GetInt("id_usuario")
	idBuzon, err := strconv.ParseInt(c.Param("id_buzon"), 10, 64)
	if err != nil {
		utils.CrearRespuesta(errors.New("No existe buzon"), nil, c, http.StatusInternalServerError)
		return
	}
	buzon := &models.Buzon{}
	buzonInicio := &models.Buzon{}
	err = models.Db.First(buzonInicio, idBuzon).Error
	if err != nil {
		utils.CrearRespuesta(errors.New("No existe buzon"), nil, c, http.StatusInternalServerError)
		return
	}
	if buzonInicio.BuzonRemitenteID != nil {
		idBuzon = int64(*buzonInicio.BuzonRemitenteID)
	}
	err = models.Db.Preload("Destinatarios").Preload("Destinatarios.Casa").Preload("Casa").Preload("Mensajes").Preload("Mensajes.Publicador").Preload("Mensajes.Casa").Preload("Mensajes.Archivos").Preload("Publicador").Preload("Archivos").First(buzon, idBuzon).Error
	for _, archivo := range buzon.Archivos {
		archivo.Url = utils.SERVIMG + archivo.Url
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe buzon"), nil, c, http.StatusInternalServerError)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
		return
	}
	if rol == "residente" {
		mensajes := []*models.Buzon{}
		for _, msj := range buzon.Mensajes {
			if msj.IsAdmin || msj.PublicadorID == uint(idUsuario) {
				mensajes = append(mensajes, msj)
			}
		}
		buzon.Mensajes = mensajes
	}

	for _, msj := range buzon.Mensajes {
		for _, arc := range msj.Archivos {
			arc.Url = utils.SERVIMG + arc.Url
		}
		var result int64
		err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: msj.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
			return
		}
		msj.Leido = result > 0
	}
	var result int64
	err = models.Db.Model(&models.BuzonLectura{}).Where(&models.BuzonLectura{BuzonID: buzon.ID, UsuarioID: uint(idUsuario)}).Count(&result).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
		return
	}
	buzon.Leido = result > 0
	buzonLectura := &models.BuzonLectura{}
	err = models.Db.FirstOrCreate(buzonLectura, &models.BuzonLectura{UsuarioID: uint(idUsuario), BuzonID: uint(idBuzon)}).Error
	buzonesLeidos := []*models.Buzon{}
	err = models.Db.Where("buzon_remitente_id = ?", idBuzon).Find(&buzonesLeidos).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
		return

	}

	for _, buz := range buzonesLeidos {
		buzonLectura := &models.BuzonLectura{}
		err = models.Db.FirstOrCreate(buzonLectura, &models.BuzonLectura{UsuarioID: uint(idUsuario), BuzonID: buz.ID}).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener buzones"), nil, c, http.StatusInternalServerError)
			return

		}

	}
	utils.CrearRespuesta(nil, buzon, c, http.StatusOK)
}

func ResponderRespuestaBuzonPrivado(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	idEtapa := c.GetInt("id_etapa")
	idCasa := uint(c.GetInt("id_casa"))
	rol := c.GetString("rol")
	idBuzon, err := strconv.ParseUint(c.Param("id_buzon"), 10, 64)
	if err != nil {
		idBuzon = 0
	}
	buzon := &models.Buzon{}
	err = c.ShouldBindJSON(buzon)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al enviar datos"), nil, c, http.StatusBadRequest)
		return
	}
	if idCasa != 0 {
		buzon.CasaID = &idCasa
	}
	buzon.PublicadorID = uint(idUsuario)
	buzon.EtapaID = uint(idEtapa)
	buzon.IsAdmin = rol == "admin-etapa"
	buzon.Publico = false
	if idBuzon != 0 {
		tx := models.Db.Begin()
		buzonRemitente := &models.Buzon{}
		err = models.Db.Preload("Archivos").Preload("BuzonRemitente").Preload("BuzonRemitente.Archivos").Preload("BuzonRemitente.Destinatarios").First(buzonRemitente, idBuzon).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.CrearRespuesta(errors.New("No existe buzon a responder"), nil, c, http.StatusBadRequest)
				return
			}
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusOK)
			return
		}
		buzonPrincipal := &models.Buzon{Titulo: buzonRemitente.Titulo, Descripcion: buzonRemitente.Descripcion, PublicadorID: buzonRemitente.PublicadorID, Publico: false, IsAdmin: buzonRemitente.IsAdmin, EtapaID: uint(idEtapa), EsRespuesta: true, CasaID: buzonRemitente.CasaID}
		err = tx.Create(buzonPrincipal).Error
		if err != nil {
			_ = tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusOK)
			return
		}
		for _, arc := range buzonRemitente.Archivos {
			err = tx.Create(&models.BuzonArchivo{Url: arc.Url, MimeType: arc.MimeType, BuzonID: buzonPrincipal.ID}).Error
			if err != nil {
				_ = tx.Rollback()
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusOK)
				return
			}
		}

		buzon.Titulo = buzonPrincipal.Titulo
		buzon.BuzonRemitenteID = &buzonPrincipal.ID
		err = tx.Create(buzon).Error
		if err != nil {
			_ = tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusOK)
			return
		}
		err = tx.Create(&models.BuzonDestinatario{CasaID: *buzonRemitente.CasaID, BuzonID: buzonPrincipal.ID, Rol: "admin-etapa"}).Error
		err = tx.Create(&models.BuzonDestinatario{CasaID: *buzonRemitente.CasaID, BuzonID: buzon.ID, Rol: "admin-etapa"}).Error
		if err != nil {
			_ = tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusOK)
			return
		}
		_ = tx.Commit()
		utils.CrearRespuesta(nil, "Respuesta creada", c, http.StatusAccepted)
	} else {
		utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusInternalServerError)
		return
	}
}

func CreateBuzon(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	idEtapa := c.GetInt("id_etapa")
	idCasa := uint(c.GetInt("id_casa"))
	rol := c.GetString("rol")
	idBuzon, err := strconv.ParseUint(c.Param("id_buzon"), 10, 64)
	if err != nil {
		idBuzon = 0
	}
	buzon := &models.Buzon{}
	err = c.ShouldBindJSON(buzon)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al enviar datos"), nil, c, http.StatusBadRequest)
		return
	}
	if idCasa != 0 {
		buzon.CasaID = &idCasa
	}
	buzon.PublicadorID = uint(idUsuario)
	buzon.EtapaID = uint(idEtapa)
	buzon.IsAdmin = rol == "admin-etapa"

	idUBuzon := uint(idBuzon)
	destinatarios := []*models.BuzonDestinatario{}
	if idBuzon != 0 {
		buzon.BuzonRemitenteID = &idUBuzon
		buzonRemitente := &models.Buzon{}
		err = models.Db.Preload("Destinatarios").First(buzonRemitente, idBuzon).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.CrearRespuesta(errors.New("No existe buzon a responder"), nil, c, http.StatusBadRequest)
				return
			}
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear buzon"), nil, c, http.StatusOK)
			return
		}
		destinatarios = buzonRemitente.Destinatarios
		buzon.Titulo = buzonRemitente.Titulo
	}
	if buzon.Destinatarios != nil {
		for _, dest := range buzon.Destinatarios {
			dest.Rol = "residente"
		}
	}
	if len(destinatarios) > 0 {

		err = models.Db.Omit("BuzonRemitente").Omit("Destinatarios").Create(buzon).Error
		for _, dest := range destinatarios {
			err = models.Db.Create(&models.BuzonDestinatario{CasaID: dest.CasaID, Rol: "residente", BuzonID: buzon.ID}).Error

		}
	} else {

		err = models.Db.Omit("BuzonRemitente").Create(buzon).Error
	}
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al crear publicacion"), nil, c, http.StatusBadRequest)
		return
	}

	utils.CrearRespuesta(nil, buzon.ID, c, http.StatusOK)
}

func CreateArchivosBuzon(c *gin.Context) {
	idBuzon, err := strconv.ParseInt(c.Param("id_buzon"), 10, 64)
	form, err := c.MultipartForm()
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener archivos"), nil, c, http.StatusBadRequest)
		return
	}
	files := form.File["archivos[]"]
	if len(files) > 10 {
		utils.CrearRespuesta(errors.New("Maximo se puede asociar 10 archivos."), nil, c, http.StatusBadRequest)
		return
	}
	num := 0

	tx := models.Db.Begin()

	for _, file := range files {

		v := fmt.Sprintf("%d", num)
		tiempo := fmt.Sprintf("%d", time.Now().Unix())
		idBu := fmt.Sprintf("%d", idBuzon)
		fileArr := strings.Split(file.Filename, ".")
		extension := fileArr[len(fileArr)-1]
		nombre := "public/img/buzon/" + tiempo + v + idBu + "." + extension
		err = c.SaveUploadedFile(file, nombre)
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al subir archivos adjuntos"), nil, c, http.StatusBadRequest)
			return
		}

		err = tx.Create(&models.BuzonArchivo{Url: nombre, BuzonID: uint(idBuzon), MimeType: extension}).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al guardar publicacion"), nil, c, http.StatusInternalServerError)
			return
		}
		num++
	}
	tx.Commit()
	utils.CrearRespuesta(nil, "Mensaje enviado con éxito", c, http.StatusOK)
}

func CreateBuzonResidente(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	idEtapa := c.GetInt("id_etapa")
	idBuzon, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		idBuzon = 0
	}
	buzon := &models.Buzon{}
	err = c.ShouldBindJSON(buzon)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al enviar datos"), nil, c, http.StatusBadRequest)
		return
	}
	buzon.PublicadorID = uint(idUsuario)
	buzon.EtapaID = uint(idEtapa)
	buzon.Publico = true
	idUBuzon := uint(idBuzon)
	if idBuzon != 0 {
		buzon.BuzonRemitenteID = &idUBuzon
	}
	err = models.Db.Create(buzon).Error
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al crear publicacion"), nil, c, http.StatusBadRequest)
		return
	}
	utils.CrearRespuesta(nil, buzon.ID, c, http.StatusOK)
}
