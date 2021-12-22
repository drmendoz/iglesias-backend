package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/auth"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateFiel(c *gin.Context) {
	res := &models.Fiel{}
	err := c.ShouldBindJSON(res)
	if res.Usuario.Usuario == "" {
		utils.CrearRespuesta(errors.New("Por favor Ingrese usuario y/o Contrasena"), nil, c, http.StatusBadRequest)
		return

	}
	if err != nil || res.Usuario.Usuario == "" {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	res.Usuario.Contrasena, _ = auth.GenerarCodigoTemporal(6)
	resComp := &models.Fiel{}
	err = models.Db.Where("Usuario.usuario = ?", res.Usuario.Usuario).Joins("Usuario").Joins("Parroquia").First(&resComp).Error
	if resComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un usuario con ese nombre de usuario"), nil, c, http.StatusNotAcceptable)
		return
	}
	err = models.Db.Where("Usuario.correo = ?", res.Usuario.Correo).Joins("Usuario").Joins("Parroquia").First(&resComp).Error
	if resComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un usuario con ese correo"), nil, c, http.StatusNotAcceptable)
		return
	}
	if errors.Is(gorm.ErrRecordNotFound, err) {
		res.ContraHash = res.Usuario.Contrasena
		clave := auth.HashPassword(res.Usuario.Contrasena)
		res.Usuario.Contrasena = clave
		tx := models.Db.Begin()
		err = tx.Create(res).Error

		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear fiel"), nil, c, http.StatusInternalServerError)
			return
		}

		if res.Usuario.Imagen == "" {
			res.Usuario.Imagen = utils.DefaultUser
		} else {
			res.Usuario, err = UploadImagePerfil(res.Usuario, tx)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear fiel"), nil, c, http.StatusInternalServerError)
				return
			}
			res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
		}
		err := mail.EnviarCambioContrasena(*res)
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al enviar mensaje con contrasena temporal"), nil, c, http.StatusInternalServerError)
			return
		}
		tx.Commit()
		utils.CrearRespuesta(nil, res, c, http.StatusCreated)

		return
	}

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear fiel"), nil, c, http.StatusInternalServerError)
		return
	}
}

func GetFieles(c *gin.Context) {
	idParroquia := c.GetInt("id_parroquia")
	fieles := []*models.Fiel{}
	var err error
	if idParroquia != 0 {
		err = models.Db.Where("parroquia_id = ?", idParroquia).Omit("Usuario.Contrasena").Joins("Usuario").Joins("Parroquia").Find(&fieles).Error
	} else {

		err = models.Db.Omit("Usuario.Contrasena").Joins("Usuario").Joins("Parroquia").Find(&fieles).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener fieles"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, usr := range fieles {
		if usr.Usuario.Imagen == "" {
			usr.Usuario.Imagen = utils.DefaultUser
		} else {
			if !strings.HasPrefix(usr.Usuario.Imagen, "https://") {
				usr.Usuario.Imagen = utils.SERVIMG + usr.Usuario.Imagen
			}
		}
	}
	utils.CrearRespuesta(err, fieles, c, http.StatusOK)
}

func UpdateFiel(c *gin.Context) {

	res := &models.Fiel{}

	err := c.ShouldBindJSON(res)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	ui, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	res.ID = uint(ui)
	adComp := &models.Fiel{}
	err = models.Db.Joins("Usuario").First(&adComp, res.ID).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			utils.CrearRespuesta(errors.New("No existe Fiel"), nil, c, http.StatusNotFound)
			return
		} else {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear fiel"), nil, c, http.StatusInternalServerError)
			return
		}

	}

	tx := models.Db.Begin()
	if res.Usuario != nil {
		if res.Usuario.Contrasena != "" {

			res.Usuario.Contrasena = auth.HashPassword(res.Usuario.Contrasena)
		}
		if res.Usuario.Imagen != "" {
			res.Usuario, err = UploadImagePerfil(res.Usuario, tx)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al actualizar fiel"), nil, c, http.StatusInternalServerError)
				return
			}
			res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
		}
		err = tx.Where("id = ?", adComp.Usuario.ID).Omit("contrasena").Updates(res.Usuario).Error
	}

	err = tx.Omit("Usuario").Updates(res).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar fiel"), nil, c, http.StatusInternalServerError)
		return
	}
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar fiel"), nil, c, http.StatusInternalServerError)
		return
	}
	err = tx.Model(&models.Fiel{}).Where("id = ?", res.ID).Updates(map[string]interface{}{
		"confirmacion": res.Confirmacion,
	}).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar fiel"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(nil, "Fiel actualizado exitosamente", c, http.StatusOK)
}

func UpdateTokenNotificacion(c *gin.Context) {
	idFiel := c.GetInt("id_fiel")
	res := &models.Fiel{}
	err := c.ShouldBindJSON(res)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al actualizar token"), nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Select("token_notificacion").Where("id = ?", idFiel).Updates(res).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar token"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, "Token actualizado con exito", c, http.StatusOK)
}

func GetFielPorId(c *gin.Context) {
	res := &models.Fiel{}
	id := c.Param("id")
	err := models.Db.Where("fiel.id = ?", id).Omit("usuarios.contrasena").Joins("Usuario").First(res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Fiel no encontrado"), nil, c, http.StatusNotFound)
			return
		}

		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener fiel"), nil, c, http.StatusInternalServerError)
		return
	}
	if res.Usuario.Imagen == "" {
		res.Usuario.Imagen = utils.DefaultUser
	} else {
		if !strings.HasPrefix(res.Usuario.Imagen, "https://") {
			res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
		}
	}

	utils.CrearRespuesta(nil, res, c, http.StatusOK)
}

func DeleteFiel(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Fiel{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar Fiel"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Fiel borrado exitosamente", c, http.StatusOK)
}

func CambiarContrasenaFiel(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	usuario := &models.Usuario{}
	err := c.ShouldBindJSON(usuario)
	if err != nil {
		utils.CrearRespuesta(errors.New("Solicitud con parametros invalidos"), nil, c, http.StatusBadRequest)
		return
	}
	usrTmp := &models.Usuario{}
	err = models.Db.Select("contrasena").First(usrTmp, idUsuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe usuario"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al cambiar contraseña"), nil, c, http.StatusInternalServerError)
		return
	}
	usuario.ViejaContrasena = auth.HashPassword(usuario.ViejaContrasena)
	if usuario.ViejaContrasena != usrTmp.Contrasena {
		utils.CrearRespuesta(errors.New("Contraseña ingresada incorrecta"), nil, c, http.StatusBadRequest)
		return
	}
	usuario.Contrasena = auth.HashPassword(usuario.Contrasena)
	err = models.Db.Select("contrasena").Where("id = ?", idUsuario).Updates(usuario).Error
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al cambiar contraseña"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Contraseña cambiada exitosamente", c, http.StatusOK)
}

func GetUFielesCount(c *gin.Context) {
	var res int64
	err := models.Db.Model(&models.Fiel{}).Count(&res).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, res, c, http.StatusOK)
}

func GetInformacionPerfil(c *gin.Context) {
	idFiel := c.GetInt("id_fiel")
	res := &models.Fiel{}
	err := models.Db.Joins("Usuario").Preload("Parroquia").Preload("Parroquia.Etapa").Preload("Parroquia.Etapa.Urbanizacion").First(res, idFiel).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error interno del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	res.Usuario.Cedula = &res.Cedula
	if res.Usuario.Imagen != "" {
		res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
	} else {
		res.Usuario.Imagen = utils.DefaultUser
	}

	utils.CrearRespuesta(nil, res, c, http.StatusOK)

}

func EditarImagenPerfilFiel(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	usuario := &models.Usuario{}
	err := c.ShouldBindJSON(usuario)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en formato de envio"), nil, c, http.StatusBadRequest)
		return
	}
	usuarioid := fmt.Sprintf("%d", idUsuario)
	if usuario.Imagen != "" {
		usuario.Imagen, err = img.FromBase64ToImage(usuario.Imagen, "usuarios/"+time.Now().Format(time.RFC3339)+usuarioid, true)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(err, nil, c, http.StatusInternalServerError)
			return
		}
		err = models.Db.Where("id = ?", idUsuario).Select("imagen").Updates(usuario).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al cambiar imagen"), nil, c, http.StatusInternalServerError)
			return
		}
		usuario.Imagen = utils.SERVIMG + usuario.Imagen
	} else {
		utils.CrearRespuesta(nil, "No se envio imagen", c, http.StatusAccepted)
		return
	}
	utils.CrearRespuesta(nil, usuario.Imagen, c, http.StatusAccepted)
}
