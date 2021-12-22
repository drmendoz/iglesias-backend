package controllers

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/auth"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"

	"github.com/drmendoz/iglesias-backend/utils/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EnviarCodigoTemporal(c *gin.Context) {
	rol := c.Param("rol")
	usuarioTemp := &models.Usuario{}
	err := c.ShouldBindJSON(usuarioTemp)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error de parametros de solicitud"), nil, c, http.StatusBadRequest)
		return
	}
	usuario := &models.Usuario{}
	switch rol {
	case "admin-master":
		admin := &models.AdminMaster{}
		err = models.Db.Where("Usuario.usuario= ?", usuarioTemp.Usuario).Joins("Usuario").First(admin).Error

		usuario = admin.Usuario
	case "fiel":
		res := &models.Fiel{}
		err = models.Db.Where("Usuario.usuario= ?", usuarioTemp.Usuario).Joins("Usuario").First(res).Error

		usuario = res.Usuario
	case "admin-parroquia":
		admin := &models.AdminParroquia{}
		err = models.Db.Where("Usuario.usuario= ?", usuarioTemp.Usuario).Joins("Usuario").First(admin).Error

		usuario = admin.Usuario
	default:
		utils.CrearRespuesta(errors.New("No existe rol"), nil, c, http.StatusBadRequest)
		return
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe usuario"), nil, c, http.StatusBadRequest)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al recuperar cuenta"), nil, c, http.StatusInternalServerError)
		return
	}
	usuario.CodigoTemporal, _ = auth.GenerarCodigoTemporal(6)

	err = models.Db.Model(&models.Usuario{}).Select("codigo_temporal").Where("id = ?", usuario.ID).Updates(usuario).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al recuperar cuenta"), nil, c, http.StatusInternalServerError)
		return
	}
	err = mail.EnviarCorreoRecover(*usuario)
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al enviar correo. Por favor comuniquese con soporte"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Se envio un codigo temporal al correo electronico "+usuario.Correo, c, http.StatusOK)
}

func CambioDeContrasena(c *gin.Context) {
	rol := c.Param("rol")
	recover := &auth.NuevaContrasena{}
	err := c.ShouldBindJSON(recover)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error en parametros de solicitud"), nil, c, http.StatusBadRequest)
		return
	}
	usuario := &models.Usuario{}
	switch rol {
	case "admin-master":
		admin := &models.AdminMaster{}
		err = models.Db.Where("Usuario.usuario= ?", recover.Usuario).Joins("Usuario").First(admin).Error
		usuario = admin.Usuario
	case "fiel":
		res := &models.Fiel{}
		err = models.Db.Where("Usuario.usuario= ?", recover.Usuario).Joins("Usuario").First(res).Error
		usuario = res.Usuario
	case "admin-parroquia":
		admin := &models.AdminParroquia{}
		err = models.Db.Where("Usuario.usuario= ?", recover.Usuario).Joins("Usuario").First(admin).Error
		usuario = admin.Usuario
	default:
		utils.CrearRespuesta(errors.New("No existe rol"), nil, c, http.StatusBadRequest)
		return
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe usuario"), nil, c, http.StatusBadRequest)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al recuperar cuenta"), nil, c, http.StatusInternalServerError)
		return
	}
	if usuario.CodigoTemporal != recover.CodigoTemporal {
		utils.CrearRespuesta(errors.New("Codigo temporal ingresado incorrecto"), nil, c, http.StatusBadRequest)
		return
	}
	usuario.Contrasena = auth.HashPassword(recover.Contrasena)
	usuario.CodigoTemporal = ""
	tx := models.Db.Begin()
	err = tx.Model(&models.Usuario{}).Select("contrasena", "codigo_temporal").Where("id = ?", usuario.ID).Updates(usuario).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al cambiar contrasena"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(nil, "Contrasena cambiada con exito", c, http.StatusOK)
}
