package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/drmendoz/iglesias-backend/auth"
	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAdministradores(c *gin.Context) {
	administradores := []*models.AdminMaster{}
	err := models.Db.Joins("Usuario").Order("Usuario.Apellido ASC").Omit("Usuario.Contrasena").Preload("Permisos").Find(&administradores).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administadores"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, adm := range administradores {
		if adm.Usuario.Imagen == "" {
			adm.Usuario.Imagen = utils.DefaultUser
		} else {
			if !strings.HasPrefix(adm.Usuario.Imagen, "https://") {
				adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
			}
		}

	}
	utils.CrearRespuesta(err, administradores, c, http.StatusOK)
}

func CreateAdministrador(c *gin.Context) {
	adm := &models.AdminMaster{}
	err := c.ShouldBindJSON(adm)
	rol := c.GetString("rol")
	isMaster := rol == "master"

	if err != nil || adm.Usuario.Usuario == "" {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	adComp := &models.AdminMaster{}
	err = models.Db.Where("Usuario.usuario = ?", adm.Usuario.Usuario).Joins("Usuario").First(&adComp).Error
	fmt.Print(err)
	if errors.Is(gorm.ErrRecordNotFound, err) {
		adm.Usuario.Contrasena, _ = auth.GenerarCodigoTemporal(6)
		clave := auth.HashPassword(adm.Usuario.Contrasena)
		adm.ContraHash = adm.Usuario.Contrasena
		adm.Usuario.Contrasena = clave
		tx := models.Db.Begin()
		if isMaster {
			adm.Permisos = models.AdminMasterPermiso{
				Autorizado:    true,
				Urbanizacion:  true,
				Etapa:         true,
				Administrador: true,
				Modulo:        true,
				Categoria:     true,
				Publicidad:    true,
				Facturacion:   true,
				Fiel:          true,
				Usuario:       true,
			}
		}

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
		err = mail.EnviarCambioContrasenaMaster(*adm)
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al enviar mensaje con contrasena temporal"), nil, c, http.StatusInternalServerError)
			return
		}
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

func UpdateAdministrador(c *gin.Context) {

	adm := &models.AdminMaster{}

	err := c.ShouldBindJSON(adm)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	ui, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adm.ID = uint(ui)
	adComp := &models.AdminMaster{}
	err = models.Db.Where("Usuario.usuario = ?", adm.Usuario.Usuario).Joins("Usuario").First(&adComp).Error
	if errors.Is(gorm.ErrRecordNotFound, err) || adm.ID == adComp.ID {
		adm.Usuario.ID = adComp.UsuarioID
		adm.Usuario.Contrasena = auth.HashPassword(adm.Usuario.Contrasena)
		tx := models.Db.Begin()

		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al editar administrador"), nil, c, http.StatusInternalServerError)
			return
		}
		if img.IsBase64(adm.Usuario.Imagen) {
			adm.Usuario, err = UploadImagePerfil(adm.Usuario, tx)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al decodificar imagen"), nil, c, http.StatusInternalServerError)
				return
			}
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
		err = tx.Omit("Usuario").Updates(adm).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
			return
		}
		if adm.Usuario != nil {
			err = tx.Where("id = ?", adm.UsuarioID).Updates(adm.Usuario).Error
			if err != nil {
				tx.Rollback()
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
				return
			}

		}
		if adm.Usuario.Imagen == "" {
			adm.Usuario.Imagen = utils.DefaultUser
		}
		tx.Commit()
		utils.CrearRespuesta(nil, adm, c, http.StatusOK)
		return
	}
	// if adComp.ID != 0 {
	// 	utils.CrearRespuesta(errors.New("Ya existe un administrador con ese correo"), nil, c, http.StatusNotAcceptable)
	// 	return
	// }

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al editar administrador"), nil, c, http.StatusInternalServerError)
		return
	}
}

func GetAdministradorPorId(c *gin.Context) {
	adm := &models.AdminMaster{}
	id := c.Param("id")
	err := models.Db.Where("admin_master.id = ?", id).Omit("usuarios.contrasena").Joins("Usuario").First(adm).Error
	if adm.Usuario.Imagen == "" {
		adm.Usuario.Imagen = utils.DefaultUser
	} else {
		adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Administrador no encontrado"), nil, c, http.StatusNotFound)
			return
		}

		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administrador"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, adm, c, http.StatusOK)
}

func DeleteAdministrador(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.AdminMaster{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar administrador"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Administrador borrado exitosamente", c, http.StatusOK)
}
