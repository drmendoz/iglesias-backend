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
	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

type AutorizacionesResponse struct {
	Fijas      []*models.Autorizacion `json:"fijas"`
	Temporales []*models.Autorizacion `json:"temporales"`
}

func GetAutorizacionesAdmin(c *gin.Context) {
	casa := c.Query("id_casa")
	tipo := c.Query("tipo")
	estado := c.Query("estado")
	var idCasa int
	var err error
	if casa == "" {
		idCasa = 0
	} else {
		idCasa, err = strconv.Atoi(casa)
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener autorizaciones"), nil, c, http.StatusInternalServerError)
		return
	}
	autorizaciones := []*models.Autorizacion{}

	err = models.Db.Limit(100).Where(&models.Autorizacion{CasaID: uint(idCasa), Tipo: tipo, Estado: estado}).Joins("Casa").Order("created_at desc").Find(&autorizaciones).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener autorizaciones"), nil, c, http.StatusInternalServerError)
		return
	}

	for _, autorizacion := range autorizaciones {
		if autorizacion.Imagen != "" {
			if !strings.HasPrefix(autorizacion.Imagen, "https://") {
				autorizacion.Imagen = utils.SERVIMG + autorizacion.Imagen
			}
		}
		autorizacion.TipoUsuario = strcase.ToCamel(strings.ToLower(autorizacion.TipoUsuario))
	}

	utils.CrearRespuesta(err, autorizaciones, c, http.StatusOK)
}

func GetAutorizaciones(c *gin.Context) {
	idCasa := uint(c.GetInt("id_casa"))

	autorizaciones := &AutorizacionesResponse{}
	autorizaciones.Fijas = []*models.Autorizacion{}
	autorizaciones.Temporales = []*models.Autorizacion{}

	err := models.Db.Limit(100).Where(&models.Autorizacion{CasaID: idCasa, Tipo: "FIJA"}).Order("created_at desc").Find(&autorizaciones.Fijas).Error
	err = models.Db.Limit(100).Where(&models.Autorizacion{CasaID: idCasa, Tipo: "TEMPORAL"}).Order("created_at desc").Find(&autorizaciones.Temporales).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener autorizaciones"), nil, c, http.StatusInternalServerError)
		return
	}

	for _, autorizacion := range autorizaciones.Fijas {
		if autorizacion.Imagen != "" {
			if !strings.HasPrefix(autorizacion.Imagen, "https://") {
				autorizacion.Imagen = utils.SERVIMG + autorizacion.Imagen
			}
		}
		autorizacion.TipoUsuario = strcase.ToCamel(strings.ToLower(autorizacion.TipoUsuario))

	}
	for _, autorizacion := range autorizaciones.Temporales {
		if !strings.HasPrefix(autorizacion.Imagen, "https://") {
			autorizacion.Imagen = utils.SERVIMG + autorizacion.Imagen
		}
		autorizacion.TipoUsuario = strcase.ToCamel(strings.ToLower(autorizacion.TipoUsuario))
	}

	utils.CrearRespuesta(err, autorizaciones, c, http.StatusOK)
}

func ValidarAutorizacion(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_etapa"))
	pin := c.Query("pin")

	autorizacion := &models.Autorizacion{}

	err := models.Db.Where(&models.Autorizacion{ParroquiaID: idParroquia, Estado: "PENDIENTE", Pin: pin}).Joins("Casa").First(&autorizacion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Autorización no encontrada"), nil, c, http.StatusNotFound)
		return
	}
	if autorizacion.Imagen != "" {
		autorizacion.Imagen = utils.SERVIMG + autorizacion.Imagen
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	autorizacion.Manzana = autorizacion.Casa.Manzana
	autorizacion.Villa = autorizacion.Casa.Villa
	utils.CrearRespuesta(err, autorizacion, c, http.StatusOK)
}

func GetAutorizacionPorId(c *gin.Context) {
	autorizacion := &models.Autorizacion{}
	id := c.Param("id")
	err := models.Db.Joins("Publicador").Joins("Usuario").First(autorizacion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Autorizacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, autorizacion, c, http.StatusOK)
}

type Respuesta struct {
	Respuesta string `json:"respuesta"`
	Pin       string `json:"pin"`
}

func CreateAutorizacion(c *gin.Context) {
	idUsuario := c.GetInt("id_usuario")
	idParroquia := c.GetInt("id_etapa")
	idCasa := c.GetInt("id_casa")
	if idUsuario == 0 || idParroquia == 0 {
		utils.CrearRespuesta(errors.New("Error al crear autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	autorizacion := &models.Autorizacion{}
	err := c.ShouldBindJSON(autorizacion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}

	autorizacion.PublicadorID = uint(idUsuario)
	autorizacion.ParroquiaID = uint(idParroquia)
	autorizacion.CasaID = uint(idCasa)
	tx := models.Db.Begin()
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	if autorizacion.Tipo == "TEMPORAL" {
		autorizacion.Pin, _ = auth.GenerarCodigoTemporal(4)
	}
	if autorizacion.Imagen != "" {
		if img.IsBase64(autorizacion.Imagen) {
			autorizacion.Imagen, err = img.FromBase64ToImage(autorizacion.Imagen, "autorizaciones/"+time.Now().Format(time.RFC3339), false)
		}
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear autorizacion "), nil, c, http.StatusInternalServerError)

			return
		}
	}
	if autorizacion.Imagen2 != "" {
		autorizacion.Imagen2, err = img.FromBase64ToImage(autorizacion.Imagen2, "autorizaciones/"+time.Now().Format(time.RFC3339), false)
		if err != nil {
			_ = c.Error(err)
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Error al crear autorizacion "), nil, c, http.StatusInternalServerError)

			return
		}
	}
	err = tx.Omit("usuario_id").Create(autorizacion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(err, "Error al crear autorización", c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	respuesta := &Respuesta{}
	respuesta.Respuesta = "Autorización creada correctamente"
	respuesta.Pin = autorizacion.Pin
	utils.CrearRespuesta(err, respuesta, c, http.StatusCreated)
}

func UpdateAutorizacion(c *gin.Context) {
	autorizacion := &models.Autorizacion{}

	err := c.ShouldBindJSON(autorizacion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id := c.Param("id")
	tx := models.Db.Begin()
	autorizacion.UpdatedAt = time.Now()

	err = tx.Omit("imagen").Where("id = ?", id).Updates(autorizacion).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	tx.Commit()
	utils.CrearRespuesta(err, "Autorizacion actualizada correctamente", c, http.StatusOK)
}

func DeleteAutorizacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Autorizacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Autorizacion eliminada exitosamente", c, http.StatusOK)
}

func DeleteAutorizaciones(c *gin.Context) {
	err := models.Db.Where("id IS NOT NULL").Delete(&models.Autorizacion{}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar autorizacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Autorizacion eliminada exitosamente", c, http.StatusOK)
}
