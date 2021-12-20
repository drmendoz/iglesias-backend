package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
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

func CreateResidente(c *gin.Context) {
	res := &models.Residente{}
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
	fechaActual := time.Now()
	res.VisualizacionBitacora = &fechaActual
	res.VisualizacionBuzon = &fechaActual
	res.VisualizacionEmprendimiento = &fechaActual
	res.VisualizacionGaleria = &fechaActual
	res.VisualizacionVotacion = &fechaActual
	res.VisualizacionCamara = &fechaActual
	res.VisualizacionAlicuota = &fechaActual
	res.VisualizacionAreaSocial = &fechaActual
	res.VisualizacionAdministradores = &fechaActual
	res.VisualizacionReservas = &fechaActual
	resComp := &models.Residente{}
	err = models.Db.Where("Usuario.usuario = ?", res.Usuario.Usuario).Joins("Usuario").Joins("Casa").First(&resComp).Error
	if resComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un residente con ese usuario"), nil, c, http.StatusNotAcceptable)
		return
	}
	err = models.Db.Where("Usuario.correo = ?", res.Usuario.Correo).Joins("Usuario").Joins("Casa").First(&resComp).Error
	if resComp.ID != 0 {
		utils.CrearRespuesta(errors.New("Ya existe un residente con ese correo"), nil, c, http.StatusNotAcceptable)
		return
	}
	if errors.Is(gorm.ErrRecordNotFound, err) {
		if res.Pdf != "" {
			uri := strings.Split(res.Pdf, ";")[0]
			if uri == "data:application/pdf" {
				nombre := fmt.Sprintf("admin-etapa-%d.pdf", time.Now().Unix())
				base64 := strings.Split(res.Pdf, ",")[1]
				err = utils.SubirPdf(nombre, base64)
				if err != nil {
					_ = c.Error(err)
					utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
					return
				}
				res.Pdf = nombre
			} else {
				res.Pdf = ""
			}
		} else {
			res.Pdf = ""
		}
		res.ContraHash = res.Usuario.Contrasena
		clave := auth.HashPassword(res.Usuario.Contrasena)
		res.Usuario.Contrasena = clave
		tx := models.Db.Begin()
		err = tx.Create(res).Error

		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear residente"), nil, c, http.StatusInternalServerError)
			return
		}

		if res.Usuario.Imagen == "" {
			res.Usuario.Imagen = utils.DefaultUser
		} else {
			res.Usuario, err = UploadImagePerfil(res.Usuario, tx)
			if err != nil {
				_ = c.Error(err)
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Error al crear residente"), nil, c, http.StatusInternalServerError)
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
		res.Casa = nil
		tx.Commit()
		utils.CrearRespuesta(nil, res, c, http.StatusCreated)

		return
	}

	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear residente"), nil, c, http.StatusInternalServerError)
		return
	}
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func sortResidentes(residentes []*models.Residente, c *gin.Context) {
	sort.SliceStable(residentes, func(i, j int) bool {
		if isNumeric(residentes[i].Casa.Manzana) && isNumeric(residentes[j].Casa.Manzana) {
			mzI, err := strconv.Atoi(residentes[i].Casa.Manzana)
			if err != nil {
				_ = c.Error(err)
			}
			mzJ, err := strconv.Atoi(residentes[j].Casa.Manzana)
			if err != nil {
				_ = c.Error(err)
			}
			return mzI < mzJ
		} else {
			return residentes[i].Casa.Manzana < residentes[j].Casa.Manzana
		}
	})
}

func GetResidente(c *gin.Context) {
	idEtapa := c.GetInt("id_etapa")
	residentes := []*models.Residente{}
	var err error
	if idEtapa != 0 {
		err = models.Db.Where("Casa.etapa_id = ?", idEtapa).Order("Casa.Manzana DESC, Casa.Villa DESC").Omit("Usuario.Contrasena").Joins("Usuario").Joins("Casa").Find(&residentes).Error
		sortResidentes(residentes, c)
	} else {

		err = models.Db.Omit("Usuario.Contrasena").Joins("Usuario").Joins("Casa").Find(&residentes).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener residentes"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, usr := range residentes {
		if usr.Usuario.Imagen == "" {
			usr.Usuario.Imagen = utils.DefaultUser
		} else {
			if !strings.HasPrefix(usr.Usuario.Imagen, "https://") {
				usr.Usuario.Imagen = utils.SERVIMG + usr.Usuario.Imagen
			}
		}
		if usr.Pdf != "" {
			usr.Pdf = "https://api.practical.com.ec/public/pdf/" + usr.Pdf
		}

	}
	utils.CrearRespuesta(err, residentes, c, http.StatusOK)
}

func UpdateResidente(c *gin.Context) {

	res := &models.Residente{}

	err := c.ShouldBindJSON(res)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	ui, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	res.ID = uint(ui)
	adComp := &models.Residente{}
	err = models.Db.Joins("Usuario").First(&adComp, res.ID).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			utils.CrearRespuesta(errors.New("No existe Residente"), nil, c, http.StatusNotFound)
			return
		} else {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al crear residente"), nil, c, http.StatusInternalServerError)
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
				utils.CrearRespuesta(errors.New("Error al actualizar residente"), nil, c, http.StatusInternalServerError)
				return
			}
			res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
		}
		err = tx.Where("id = ?", adComp.Usuario.ID).Omit("contrasena").Updates(res.Usuario).Error
	}

	if res.Pdf != "" {
		uri := strings.Split(res.Pdf, ";")[0]
		if uri == "data:application/pdf" {
			nombre := fmt.Sprintf("admin-etapa-%d.pdf", time.Now().Unix())
			base64 := strings.Split(res.Pdf, ",")[1]
			err = utils.SubirPdf(nombre, base64)
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al crear administrador"), nil, c, http.StatusInternalServerError)
				return
			}
			res.Pdf = nombre
		} else {
			res.Pdf = ""
		}
	} else {
		res.Pdf = ""
	}
	err = tx.Omit("Usuario").Updates(res).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar residente"), nil, c, http.StatusInternalServerError)
		return
	}
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar residente"), nil, c, http.StatusInternalServerError)
		return
	}
	err = tx.Model(&models.Residente{}).Where("id = ?", res.ID).Updates(map[string]interface{}{
		"confirmacion": res.Confirmacion,
		"autorizacion": res.Autorizacion,
		"is_principal": res.IsPrincipal}).Error
	if err != nil {
		tx.Rollback()
		utils.CrearRespuesta(errors.New("Error al actualizar residente"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(nil, "Residente actualizado exitosamente", c, http.StatusOK)
}

func UpdateTokenNotificacion(c *gin.Context) {
	idResidente := c.GetInt("id_residente")
	res := &models.Residente{}
	err := c.ShouldBindJSON(res)
	if err != nil {
		utils.CrearRespuesta(errors.New("Error al actualizar token"), nil, c, http.StatusBadRequest)
		return
	}
	err = models.Db.Select("token_notificacion").Where("id = ?", idResidente).Updates(res).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar token"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(err, "Token actualizado con exito", c, http.StatusOK)
}

func GetResidentePorId(c *gin.Context) {
	res := &models.Residente{}
	id := c.Param("id")
	err := models.Db.Where("residente.id = ?", id).Omit("usuarios.contrasena").Joins("Usuario").First(res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("residente no encontrado"), nil, c, http.StatusNotFound)
			return
		}

		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener residente"), nil, c, http.StatusInternalServerError)
		return
	}
	if res.Usuario.Imagen == "" {
		res.Usuario.Imagen = utils.DefaultUser
	} else {
		if !strings.HasPrefix(res.Usuario.Imagen, "https://") {
			res.Usuario.Imagen = utils.SERVIMG + res.Usuario.Imagen
		}
	}
	if res.Pdf != "" {
		res.Pdf = "https://api.practical.com.ec/public/pdf/" + res.Pdf
	}
	utils.CrearRespuesta(nil, res, c, http.StatusOK)
}

func DeleteResidente(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Residente{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar Residente"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Residente borrado exitosamente", c, http.StatusOK)
}

func GetResidentesPorCasa(c *gin.Context) {
	casa := &models.Casa{}
	id := c.Param("id")
	err := models.Db.Preload("Residentes").Preload("Residentes.Usuario").First(casa, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Casa no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener casa"), nil, c, http.StatusInternalServerError)
		return
	}
	if casa.Imagen == "" {
		casa.Imagen = utils.DefaultEtapa
	} else {
		casa.Imagen = utils.SERVIMG + casa.Imagen
	}

	for _, residente := range casa.Residentes {
		if residente.Usuario.Imagen == "" {
			residente.Usuario.Imagen = utils.DefaultUser
		} else {
			residente.Usuario.Imagen = utils.SERVIMG + residente.Usuario.Imagen
		}
	}

	utils.CrearRespuesta(nil, casa, c, http.StatusOK)
}

func CambiarContrasenaResidente(c *gin.Context) {
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

func GetUResidentesCount(c *gin.Context) {
	var res int64
	err := models.Db.Model(&models.Residente{}).Count(&res).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error del servidor"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, res, c, http.StatusOK)
}

func GetInformacionPerfil(c *gin.Context) {
	idResidente := c.GetInt("id_residente")
	res := &models.Residente{}
	err := models.Db.Joins("Usuario").Preload("Casa").Preload("Casa.Etapa").Preload("Casa.Etapa.Urbanizacion").First(res, idResidente).Error
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

func EditarImagenPerfilResidente(c *gin.Context) {
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
