package controllers

import (
	"errors"
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

func GetAdministradoresParroquia(c *gin.Context) {
	administradores := []*models.AdminParroquia{}
	idParroquia := c.GetInt("id_parroquia")
	err := models.Db.Where(&models.AdminParroquia{ParroquiaID: uint(idParroquia)}).Omit("usuario.Contrasena").Joins("Parroquia").Joins("Usuario").Order("Usuario.Apellido ASC").Preload("Permisos").Find(&administradores).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administadores"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, adm := range administradores {
		if adm.Usuario.Imagen == "" {
			adm.Usuario.Imagen = utils.DefaultUser
		} else if !strings.HasPrefix(adm.Usuario.Imagen, "https://") {
			adm.Usuario.Imagen = "https://api.practical.com.ec/public/pdf/" + adm.Usuario.Imagen
		}
		adm.Usuario.Contrasena = ""
	}
	utils.CrearRespuesta(err, administradores, c, http.StatusOK)
}

func CreateAdministradorParroquia(c *gin.Context) {

	idParroquia := c.GetInt("id_parroquia")
	adm := &models.AdminParroquia{}
	rol := c.GetString("rol")
	isMaster := rol == "master"

	err := c.ShouldBindJSON(adm)
	if idParroquia != 0 {
		adm.ParroquiaID = uint(idParroquia)
	}
	if err != nil || adm.Usuario.Correo == "" {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	if isMaster {
		adm.Permisos = models.AdminParroquiaPermiso{Usuario: true, Horario: true,
			Emprendimiento: true, Actividad: true, Intencion: true, Musica: true,
			Ayudemos: true, Misa: true,
			Curso: true}
		adm.EsMaster = true
	}
	adComp := &models.AdminParroquia{}
	err = models.Db.Where("Usuario.correo = ?", adm.Usuario.Correo).Joins("Usuario").First(&adComp).Error

	if errors.Is(gorm.ErrRecordNotFound, err) {

		adm.Usuario.Contrasena, _ = auth.GenerarCodigoTemporal(6)
		clave := auth.HashPassword(adm.Usuario.Contrasena)
		adm.ContraHash = adm.Usuario.Contrasena
		adm.Usuario.Contrasena = clave
		tx := models.Db.Begin()
		err = tx.Create(adm).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
			return
		}

		if adm.Usuario.Imagen == "" {
			adm.Usuario.Imagen = utils.DefaultUser
		} else {
			adm.Usuario, err = UploadImagePerfil(adm.Usuario, tx)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
				return
			}
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
		_ = mail.EnviarCambioContrasenaParroquia(*adm)
		tx.Commit()
		utils.CrearRespuesta(nil, adm, c, http.StatusCreated)

		return
	}
	if adComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un administrador con ese correo"), nil, c, http.StatusNotAcceptable)
		return
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
		return
	}
}

func UpdateAdministradorParroquia(c *gin.Context) {

	adm := &models.AdminParroquia{}

	err := c.ShouldBindJSON(adm)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	ui, _ := strconv.Atoi(c.Param("id"))
	adm.ID = uint(ui)
	adComp := &models.AdminParroquia{}
	if adm.Usuario != nil {

		err = models.Db.Where("Usuario.correo = ?", adm.Usuario.Correo).Joins("Usuario").First(&adComp).Error
	}

	if adm.Usuario == nil || errors.Is(gorm.ErrRecordNotFound, err) || adm.ID == adComp.ID {
		tx := models.Db.Begin()
		err = tx.Omit("Usuario, Parroquia").Updates(adm).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
			return
		}
		if adm.Usuario != nil {
			err = tx.Omit("imagen", "contrasena, Parroquia").Where("id = ?", adm.UsuarioID).Updates(adm.Usuario).Error
			if err != nil {
				tx.Rollback()
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
				return
			}
			if img.IsBase64(adm.Usuario.Imagen) {
				img, err := img.FromBase64ToImage(adm.Usuario.Imagen, "usuarios/"+time.Now().Format(time.RFC3339), false)
				if err != nil {
					_ = c.Error(err)
					tx.Rollback()
					utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
					return
				}
				adm.Usuario.Imagen = img
				err = tx.Model(&models.Usuario{}).Where("id = ?", adComp.UsuarioID).Update("imagen", img).Error
				if err != nil {
					_ = c.Error(err)
					tx.Rollback()
					utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
					return
				}
			}
		}
		err = tx.Model(&models.AdminParroquia{}).Where("id = ?", ui).Updates(map[string]interface{}{
			"es_master": adm.EsMaster,
		}).Error
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al editar administrador"), nil, c, http.StatusInternalServerError)
			return
		}

		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al editar administrador"), nil, c, http.StatusInternalServerError)
			return
		}

		tx.Commit()
		utils.CrearRespuesta(nil, "Administrador actualizado", c, http.StatusOK)
		return
	}
	if adComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un administrador con ese correo"), nil, c, http.StatusNotAcceptable)
		return
	}

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al editar administrador"), nil, c, http.StatusInternalServerError)
		return
	}
}

func GetAdministradorParroquiaPorId(c *gin.Context) {
	adm := &models.AdminMaster{}
	id := c.Param("id")
	err := models.Db.Where("admin_etapa.id = ?", id).Omit("usuarios.contrasena").Joins("Usuario").Preload("Permisos").First(adm).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Administrador no encontrado"), nil, c, http.StatusNotFound)
			return
		}

		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administrador"), nil, c, http.StatusInternalServerError)
		return
	}
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser
	} else {
		adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
	}

	utils.CrearRespuesta(nil, adm, c, http.StatusOK)
}

func DeleteAdministradorParroquia(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.AdminParroquia{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("No existe administrador"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar administrador"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Administrador borrado exitosamente", c, http.StatusOK)
}
