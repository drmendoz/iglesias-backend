package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/auth"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginAdminMaster(c *gin.Context) {
	creds := &auth.Login{}
	err := c.ShouldBindJSON(creds)

	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros de Request Invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	adm := &models.AdminMaster{}
	creds.Contrasena = auth.HashPassword(creds.Contrasena)
	res := models.Db.Where("Usuario.usuario = ?", creds.Usuario).Joins("Usuario").Preload("Permisos").First(adm)
	if res.Error != nil || creds.Contrasena != adm.Usuario.Contrasena {
		utils.CrearRespuesta(errors.New("Usuario y/o contrasena incorrecta"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser
	} else {
		if !strings.HasPrefix(adm.Usuario.Imagen, "https://") {
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
	}

	adm.Token = auth.GenerarToken(adm.Usuario, "admin-master")
	adm.Usuario.Contrasena = ""
	utils.CrearRespuesta(nil, adm, c, http.StatusAccepted)
}

func LoginAdminGarita(c *gin.Context) {
	creds := &auth.Login{}
	err := c.ShouldBindJSON(creds)

	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros de Request Invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	adm := &models.AdminGarita{}
	creds.Contrasena = auth.HashPassword(creds.Contrasena)
	res := models.Db.Where("Usuario.usuario = ? ", creds.Usuario).Joins("Usuario").Joins("Etapa").Preload("Etapa.Urbanizacion").First(adm)
	if res.Error != nil || creds.Contrasena != adm.Usuario.Contrasena {
		utils.CrearRespuesta(errors.New("Usuario y/o contrasena incorrecta"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser
	} else {

		adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
	}
	if adm.Etapa == nil {
		utils.CrearRespuesta(errors.New("Su etapa ya no existe"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Etapa.Urbanizacion == nil {
		utils.CrearRespuesta(errors.New("Su urbanizacion ya no existe"), nil, c, http.StatusUnauthorized)
		return
	}

	adm.Token = auth.GenerarToken(adm.Usuario, "admin-garita")
	if adm.Etapa.Imagen == "" {
		adm.Etapa.Imagen = utils.DefaultEtapa
	} else {
		adm.Etapa.Imagen = utils.SERVIMG + adm.Etapa.Imagen
	}
	if adm.Etapa.Urbanizacion.Imagen == "" {
		adm.Etapa.Urbanizacion.Imagen = utils.DefaultEtapa
	} else {
		adm.Etapa.Urbanizacion.Imagen = utils.SERVIMG + adm.Etapa.Urbanizacion.Imagen
	}
	adm.EtapaLabel = &models.EtapaInfo{
		EtapaNombre: adm.Etapa.Nombre,
		Imagen:      adm.Etapa.Imagen,
	}
	adm.UrbanizacionLabel = &models.UrbInfo{
		Nombre: adm.Etapa.Urbanizacion.Nombre,
		Imagen: adm.Etapa.Urbanizacion.Imagen,
	}
	utils.CrearRespuesta(nil, adm, c, http.StatusAccepted)
}

func LoginAdminEtapa(c *gin.Context) {
	creds := &auth.Login{}
	err := c.ShouldBindJSON(creds)

	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros de Request Invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	adm := &models.AdminEtapa{}
	creds.Contrasena = auth.HashPassword(creds.Contrasena)
	res := models.Db.Where("Usuario.usuario = ? ", creds.Usuario).Joins("Usuario").Preload("Etapa", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "nombre", "imagen", "urbanizacion_id", "pagos_tarjeta",
			"modulo_market", "modulo_publicacion", "modulo_votacion", "modulo_area_social",
			"modulo_equipo", "modulo_historia", "modulo_bitacora", "urbanizacion", "formulario_entrada",
			"formulario_salida", "modulo_alicuota", "modulo_emprendimiento", "modulo_camaras", "modulo_directiva",
			"modulo_galeria", "modulo_horarios", "modulo_mi_registro")
	}).Preload("Etapa.Urbanizacion", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "nombre")
	}).Preload("Permisos").First(adm)
	if res.Error != nil || creds.Contrasena != adm.Usuario.Contrasena {
		utils.CrearRespuesta(errors.New("Usuario y/o contrasena incorrecta"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Etapa.Urbanizacion == nil {
		utils.CrearRespuesta(errors.New("Su urbanizacion ya no existe"), nil, c, http.StatusUnauthorized)
		return

	}
	if adm.Etapa == nil {
		utils.CrearRespuesta(errors.New("Su etapa ya no existe"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser
	} else {
		if !strings.HasPrefix(adm.Usuario.Imagen, "https://") {
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
	}
	adm.Modulos = &models.Modulos{
		ModuloMarket:               adm.Etapa.ModuloMarket,
		ModuloPublicacion:          adm.Etapa.ModuloPublicacion,
		ModuloVotacion:             adm.Etapa.ModuloVotacion,
		ModuloAreaSocial:           adm.Etapa.ModuloAreaSocial,
		ModuloEquipoAdministrativo: adm.Etapa.ModuloEquipoAdministrativo,
		ModuloHistoria:             adm.Etapa.ModuloHistoria,
		ModuloBitacora:             adm.Etapa.ModuloBitacora,
	}
	adm.NombreEtapa = &adm.Etapa.Nombre
	adm.NombreUrbanizacion = &adm.Etapa.Urbanizacion.Nombre

	adm.Etapa = nil
	adm.Token = auth.GenerarToken(adm.Usuario, "admin-etapa")
	utils.CrearRespuesta(nil, adm, c, http.StatusAccepted)
}

func LoginResidente(c *gin.Context) {
	creds := &auth.Login{}
	err := c.ShouldBindJSON(creds)

	if err != nil {
		utils.Log.Warn(err)
		utils.CrearRespuesta(errors.New("Parametros de Request Invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	adm := &models.Residente{}
	creds.Contrasena = auth.HashPassword(creds.Contrasena)
	err = models.Db.Where("Usuario.usuario = ? ", creds.Usuario).Joins("Usuario").Preload("Casa").Preload("Casa.Etapa").Preload("Casa.Etapa.Urbanizacion").First(adm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Usuario y/o contrasena incorrecta"), nil, c, http.StatusUnauthorized)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error interno del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	if creds.Contrasena != adm.Usuario.Contrasena {
		utils.CrearRespuesta(errors.New("Usuario y/o contrasena incorrecta"), nil, c, http.StatusUnauthorized)
		return

	}

	if adm.Casa == nil {
		utils.CrearRespuesta(errors.New("Su casa ha sido eliminada del sistema"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Casa.Etapa == nil {
		utils.CrearRespuesta(errors.New("Su etapa ha sido eliminada del sistema"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Casa.Etapa.Urbanizacion == nil {
		utils.CrearRespuesta(errors.New("Su urbanizacion ha sido eliminada del sistema"), nil, c, http.StatusUnauthorized)
		return
	}
	adm.Usuario.Contrasena = ""
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser

	} else {
		adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
	}
	if adm.Casa.Imagen == "" {
		adm.Casa.Imagen = utils.DefaultCasa

	} else {
		adm.Casa.Imagen = utils.SERVIMG + adm.Casa.Imagen
	}
	if adm.Casa.Etapa.Imagen == "" {
		adm.Casa.Etapa.Imagen = utils.DefaultEtapa

	} else {
		adm.Casa.Etapa.Imagen = utils.SERVIMG + adm.Casa.Etapa.Imagen
	}
	if adm.Casa.Etapa.Urbanizacion.Imagen == "" {
		adm.Casa.Etapa.Urbanizacion.Imagen = utils.DefaultUrb

	} else {
		adm.Casa.Etapa.Urbanizacion.Imagen = utils.SERVIMG + adm.Casa.Etapa.Urbanizacion.Imagen
	}
	adm.Usuario.Cedula = &adm.Cedula
	numTemporal, _ := auth.GenerarCodigoTemporal(6)
	err = models.Db.Model(&adm.Usuario).Updates(models.Usuario{RandomNumToken: numTemporal}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New(("Error al iniciar sesion")), nil, c, http.StatusInternalServerError)
		return
	}
	adm.Token = auth.GenerarToken(adm.Usuario, "residente")
	adm.TokenNotificacion = creds.TokenNotificacion
	if adm.Confirmacion {
		adm.Mensaje = "Es necesario cambiar contrasena"
		utils.CrearRespuesta(nil, adm, c, http.StatusOK)
		return
	}

	err = models.Db.Model(&models.Residente{}).Where("id = ?", adm.ID).Updates(models.Residente{SesionIniciada: true}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New(("Error al iniciar sesion")), nil, c, http.StatusInternalServerError)
		return
	}
	adm.Usuario.Contrasena = ""
	utils.CrearRespuesta(nil, adm, c, http.StatusAccepted)
}

func CambioDeContrasenaResidente(c *gin.Context) {
	recover := &auth.NuevaContrasena{}
	err := c.ShouldBindJSON(recover)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en parametros de solicitud"), nil, c, http.StatusBadRequest)
		return
	}
	usuario := &models.Usuario{}
	res := &models.Residente{}
	err = models.Db.Where("Usuario.usuario= ?", recover.Usuario).Joins("Usuario").First(res).Error
	usuario = res.Usuario
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe usuario"), nil, c, http.StatusBadRequest)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al recuperar cuenta"), nil, c, http.StatusInternalServerError)
		return
	}
	recover.CodigoTemporal = auth.HashPassword(recover.CodigoTemporal)
	// if res.Usuario.Contrasena != recover.CodigoTemporal {
	// 	utils.CrearRespuesta(errors.New("Codigo temporal ingresado incorrecto"), nil, c, http.StatusBadRequest)
	// 	return
	// }
	recover.Contrasena = auth.HashPassword(recover.Contrasena)
	tx := models.Db.Begin()

	if recover.Imagen != "" {
		usuarioID := fmt.Sprintf("%d", usuario.ID)
		recover.Imagen, err = img.FromBase64ToImage(recover.Imagen, "usuarios/"+time.Now().Format(time.RFC3339)+usuarioID, false)
		if err != nil {
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al convertir imagen"), nil, c, http.StatusInternalServerError)
			return
		}
		usuario.Imagen = recover.Imagen
		err = tx.Model(&models.Usuario{}).Where("id = ?", usuario.ID).Update("imagen", recover.Imagen).Error
		if err != nil {
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al cambiar imagen"), nil, c, http.StatusInternalServerError)
			return
		}
		res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
	}
	println("ewcover.Imagen")
	println(recover.Imagen)
	println(recover.Imagen != "")
	err = tx.Model(&models.Usuario{}).Select("contrasena").Where("id = ?", usuario.ID).Updates(models.Usuario{Contrasena: recover.Contrasena}).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cambiar contrasena"), nil, c, http.StatusInternalServerError)
		return
	}
	res.Confirmacion = false
	err = tx.Model(&models.Residente{}).Select("confirmacion").Where("id = ?", res.ID).Updates(res).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cambiar contrasena"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(nil, "Contrasena cambiada con exito", c, http.StatusOK)
}

func CerrarSesion(c *gin.Context) {
	idResidente := c.GetInt("id_residente")
	tx := models.Db.Begin()
	err := tx.Model(&models.Residente{}).Where(" id = ?", idResidente).Update("sesion_iniciada", false).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cerrar sesion"), nil, c, http.StatusInternalServerError)
		return
	}
	err = tx.Model(&models.Residente{}).Where(" id = ?", idResidente).Update("token_notificacion", "").Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cerrar sesion"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(nil, "Cierre de sesion exitoso", c, http.StatusAccepted)
}
