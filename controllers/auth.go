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
	res := models.Db.Where("Usuario.correo = ?", creds.Correo).Joins("Usuario").Preload("Permisos").First(adm)
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

func LoginAdminParroquia(c *gin.Context) {
	creds := &auth.Login{}
	err := c.ShouldBindJSON(creds)

	if err != nil {
		utils.CrearRespuesta(errors.New("Parametros de Request Invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	adm := &models.AdminParroquia{}
	creds.Contrasena = auth.HashPassword(creds.Contrasena)
	res := models.Db.Where("Usuario.correo = ? ", creds.Correo).Joins("Usuario").Preload("Parroquia").Preload("Parroquia.Modulos").Preload("Permisos").First(adm)
	if res.Error != nil || creds.Contrasena != adm.Usuario.Contrasena {
		utils.CrearRespuesta(errors.New("Usuario y/o contrasena incorrecta"), nil, c, http.StatusUnauthorized)
		return
	}

	if adm.Parroquia == nil {
		utils.CrearRespuesta(errors.New("Su parroquia ya no existe"), nil, c, http.StatusUnauthorized)
		return
	}
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser
	} else {
		if !strings.HasPrefix(adm.Usuario.Imagen, "https://") {
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
	}
	adm.Token = auth.GenerarToken(adm.Usuario, "admin-parroquia")
	utils.CrearRespuesta(nil, adm, c, http.StatusAccepted)
}

func LoginFiel(c *gin.Context) {
	creds := &auth.Login{}
	err := c.ShouldBindJSON(creds)

	if err != nil {
		utils.Log.Warn(err)
		utils.CrearRespuesta(errors.New("Parametros de Request Invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	adm := &models.Fiel{}
	creds.Contrasena = auth.HashPassword(creds.Contrasena)
	err = models.Db.Where("Usuario.correo = ? ", creds.Correo).Joins("Usuario").Preload("Parroquia").First(adm).Error
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

	adm.Usuario.Contrasena = ""
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser

	} else {
		adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
	}
	if adm.Parroquia.Imagen == "" {
		adm.Parroquia.Imagen = utils.DefaultCasa

	} else {
		adm.Parroquia.Imagen = utils.SERVIMG + adm.Parroquia.Imagen
	}

	adm.Usuario.Cedula = &adm.Cedula
	numTemporal, _ := auth.GenerarCodigoTemporal(6)
	err = models.Db.Model(&adm.Usuario).Updates(models.Usuario{RandomNumToken: numTemporal}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New(("Error al iniciar sesion")), nil, c, http.StatusInternalServerError)
		return
	}
	adm.Token = auth.GenerarToken(adm.Usuario, "fiel")
	adm.TokenNotificacion = creds.TokenNotificacion
	if adm.Confirmacion {
		adm.Mensaje = "Es necesario cambiar contrasena"
		utils.CrearRespuesta(nil, adm, c, http.StatusOK)
		return
	}

	err = models.Db.Model(&models.Fiel{}).Where("id = ?", adm.ID).Updates(models.Fiel{SesionIniciada: true}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New(("Error al iniciar sesion")), nil, c, http.StatusInternalServerError)
		return
	}
	adm.Usuario.Contrasena = ""
	utils.CrearRespuesta(nil, adm, c, http.StatusAccepted)
}

func CambioDeContrasenaFiel(c *gin.Context) {
	recover := &auth.NuevaContrasena{}
	err := c.ShouldBindJSON(recover)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en parametros de solicitud"), nil, c, http.StatusBadRequest)
		return
	}
	usuario := &models.Usuario{}
	res := &models.Fiel{}
	err = models.Db.Where("Usuario.correo= ?", recover.Correo).Joins("Usuario").First(res).Error
	usuario = res.Usuario
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Correo no válido"), nil, c, http.StatusBadRequest)
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
	err = tx.Model(&models.Usuario{}).Select("contrasena").Where("id = ?", usuario.ID).Updates(models.Usuario{Contrasena: recover.Contrasena}).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cambiar contrasena"), nil, c, http.StatusInternalServerError)
		return
	}
	res.Confirmacion = false
	err = tx.Model(&models.Fiel{}).Select("confirmacion").Where("id = ?", res.ID).Updates(res).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cambiar contrasena"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(nil, "Contrasena cambiada con exito", c, http.StatusOK)
}

func CerrarSesion(c *gin.Context) {
	idFiel := c.GetInt("id_fiel")
	tx := models.Db.Begin()
	err := tx.Model(&models.Fiel{}).Where(" id = ?", idFiel).Update("sesion_iniciada", false).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cerrar sesion"), nil, c, http.StatusInternalServerError)
		return
	}
	err = tx.Model(&models.Fiel{}).Where(" id = ?", idFiel).Update("token_notificacion", "").Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cerrar sesion"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(nil, "Cierre de sesion exitoso", c, http.StatusAccepted)
}
