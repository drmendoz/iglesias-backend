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
	"github.com/drmendoz/iglesias-backend/utils/mail"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAdministradoresGarita(c *gin.Context) {
	administradores := []*models.AdminGarita{}
	idEtapa := c.GetInt("id_etapa")
	var err error
	if idEtapa != 0 {
		err = models.Db.Omit("usuarios.Contrasena").Where("etapa_id = ?", idEtapa).Joins("Usuario").Find(&administradores).Error
	} else {

		err = models.Db.Omit("usuarios.Contrasena").Joins("Usuario").Find(&administradores).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener administadores"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, adm := range administradores {
		if adm.Usuario.Imagen == "" {
			adm.Usuario.Imagen = utils.DefaultUser
		} else {
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
		if adm.Pdf != "" {
			adm.Pdf = "https://api.practical.com.ec/public/pdf/" + adm.Pdf
		}

	}
	utils.CrearRespuesta(err, administradores, c, http.StatusOK)
}

func CreateAdministradorGarita(c *gin.Context) {
	idEtapa := uint(c.GetInt("id_etapa"))
	adm := &models.AdminGarita{}
	err := c.ShouldBindJSON(adm)
	if adm.Usuario.Usuario == "" {
		utils.CrearRespuesta(errors.New("Por favor Ingrese usuario y/o Contrasena"), nil, c, http.StatusBadRequest)
		return
	}
	adm.EtapaID = idEtapa
	if err != nil || adm.Usuario.Usuario == "" {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	adComp := &models.AdminGarita{}
	err = models.Db.Where("Usuario.usuario = ?", adm.Usuario.Usuario).Joins("Usuario").First(&adComp).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		if adm.Pdf != "" {
			uri := strings.Split(adm.Pdf, ";")[0]
			if uri == "data:application/pdf" {
				nombre := fmt.Sprintf("admin-garita-%d.pdf", time.Now().Unix())
				base64 := strings.Split(adm.Pdf, ",")[1]
				err = utils.SubirPdf(nombre, base64)
				if err != nil {
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
					return
				}
				adm.Pdf = nombre
			} else {
				adm.Pdf = ""
			}
		} else {
			adm.Pdf = ""
		}
		adm.Usuario.Contrasena, _ = auth.GenerarCodigoTemporal(6)
		clave := auth.HashPassword(adm.Usuario.Contrasena)
		adm.ContraHash = adm.Usuario.Contrasena
		adm.Usuario.Contrasena = clave
		tx := models.Db.Begin()
		err = tx.Create(adm).Error
		if err != nil {
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
			_ = c.Error(err)
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
		err = mail.EnviarCambioContrasenaGarita(*adm)
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

func UpdateAdministradorgarita(c *gin.Context) {

	adm := &models.AdminGarita{}

	err := c.ShouldBindJSON(adm)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	ui, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	adm.ID = uint(ui)
	adComp := &models.AdminGarita{}
	err = models.Db.Where("Usuario.usuario = ?", adm.Usuario.Usuario).Joins("Usuario").First(&adComp).Error
	if errors.Is(gorm.ErrRecordNotFound, err) || adm.ID == adComp.ID {

		if adm.Pdf != "" {
			uri := strings.Split(adm.Pdf, ";")[0]
			if uri == "data:application/pdf" {
				nombre := fmt.Sprintf("admin-garita-%d.pdf", time.Now().Unix())
				base64 := strings.Split(adm.Pdf, ",")[1]
				err = utils.SubirPdf(nombre, base64)
				if err != nil {
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
					return
				}
				adm.Pdf = nombre
			} else {
				adm.Pdf = ""
			}
		} else {
			adm.Pdf = ""
		}
		adm.Usuario.Contrasena = auth.HashPassword(adm.Usuario.Contrasena)
		tx := models.Db.Begin()
		if adm.Usuario.Imagen != "" {
			adm.Usuario, err = UploadImagePerfil(adm.Usuario, tx)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
				return
			}
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}

		err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Omit("id_etapa").Updates(adm).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
			return
		}
		if adm.Usuario.Imagen == "" {
			adm.Usuario.Imagen = utils.DefaultUser
		} else {
			adm.Usuario.Imagen = utils.SERVIMG + adm.Usuario.Imagen
		}
		tx.Commit()
		utils.CrearRespuesta(nil, adm, c, http.StatusOK)
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

func GetAdministradorGaritaPorId(c *gin.Context) {
	adm := &models.AdminGarita{}
	id := c.Param("id")
	err := models.Db.Where("admin_urbanizacion.id = ?", id).Omit("usuarios.contrasena").Joins("Usuario").First(adm).Error
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

func DeleteAdministradorGarita(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.AdminGarita{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar administrador"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Administrador borrado exitosamente", c, http.StatusOK)
}
