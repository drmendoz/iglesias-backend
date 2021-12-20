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

func GetAdministradoresEtapa(c *gin.Context) {
	administradores := []*models.AdminParroquia{}
	idEtapa := c.GetInt("id_etapa")
	err := models.Db.Where(&models.AdminParroquia{EtapaID: uint(idEtapa)}).Omit("usuario.Contrasena").Joins("Usuario").Order("Usuario.Apellido ASC").Preload("Permisos").Find(&administradores).Error
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
		if adm.Pdf != "" {
			adm.Pdf = utils.SERVIMG + adm.Pdf
		}

	}
	utils.CrearRespuesta(err, administradores, c, http.StatusOK)
}

func CreateAdministradorEtapa(c *gin.Context) {

	idEtapa := c.GetInt("id_etapa")
	adm := &models.AdminParroquia{}
	rol := c.GetString("rol")
	isMaster := rol == "master"

	err := c.ShouldBindJSON(adm)
	if idEtapa != 0 {
		adm.EtapaID = uint(idEtapa)
	}
	if adm.Usuario.Usuario == "" {
		utils.CrearRespuesta(errors.New("Por favor Ingrese usuario y/o Contrasena"), nil, c, http.StatusBadRequest)
		return
	}
	if err != nil || adm.Usuario.Usuario == "" {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	if isMaster {
		adm.Permisos = models.AdminParroquiaPermiso{Alicuota: true, AreaSocial: true,
			Emprendimiento: true, Casa: true, Usuario: true, Seguridad: true,
			Ingreso: true, Voto: true, Directiva: true, Camara: true, Reserva: true,
			ExpresoEscolar: true, Buzon: true}
		adm.EsMaster = true
	}
	adComp := &models.AdminParroquia{}
	err = models.Db.Where("Usuario.usuario = ?", adm.Usuario.Usuario).Joins("Usuario").First(&adComp).Error

	if errors.Is(gorm.ErrRecordNotFound, err) {
		if adm.Pdf != "" {
			uri := strings.Split(adm.Pdf, ";")[0]
			if uri == "data:application/pdf" {
				nombre := fmt.Sprintf("admin-etapa-%d.pdf", time.Now().Unix())
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
		err = mail.EnviarCambioContrasenaEtapa(*adm)
		tx.Commit()
		utils.CrearRespuesta(nil, adm, c, http.StatusCreated)

		return
	}
	if adComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un administrador con ese usuario"), nil, c, http.StatusNotAcceptable)
		return
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
		return
	}
}

func UpdateAdministradorEtapa(c *gin.Context) {

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

		err = models.Db.Where("Usuario.usuario = ?", adm.Usuario.Usuario).Joins("Usuario").First(&adComp).Error
	}
	if adm.Pdf != "" {
		uri := strings.Split(adm.Pdf, ";")[0]
		if uri == "data:application/pdf" {
			nombre := fmt.Sprintf("admin-etapa-%d.pdf", time.Now().Unix())
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
	if adm.Usuario == nil || errors.Is(gorm.ErrRecordNotFound, err) || adm.ID == adComp.ID {
		tx := models.Db.Begin()
		err = tx.Omit("Usuario").Updates(adm).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al actualizar administrador"), nil, c, http.StatusInternalServerError)
			return
		}
		if adm.Usuario != nil {
			err = tx.Omit("imagen").Where("id = ?", adm.UsuarioID).Updates(adm.Usuario).Error
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

func GetAdministradorEtapaPorId(c *gin.Context) {
	adm := &models.AdminMaster{}
	id := c.Param("id")
	err := models.Db.Where("admin_etapa.id = ?", id).Omit("usuarios.contrasena").Joins("Usuario").First(adm).Error

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

func DeleteAdministradorEtapa(c *gin.Context) {
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
