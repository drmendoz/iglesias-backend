package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/drmendoz/iglesias-backend/utils/notification"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPublicacions(c *gin.Context) {
	publicacions := []*models.Publicacion{}
	var err error
	idParroquia := c.GetInt("id_etapa")
	if idParroquia != 0 {
		err = models.Db.Order("created_at desc").Where("publicacion.etapa_id = ?", idParroquia).Joins("Usuario").Joins("Etapa").Preload("ImagenesPublicacion").Find(&publicacions).Error
	} else {
		err = models.Db.Find(&publicacions).Error
	}
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener publicacions"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, publicacion := range publicacions {
		for _, img := range publicacion.ImagenesPublicacion {
			if img.Imagen != "" {
				img.Imagen = utils.SERVIMG + img.Imagen
			} else {
				img.Imagen = utils.DefaultNoticia
			}
		}
	}
	utils.CrearRespuesta(nil, publicacions, c, http.StatusOK)
}

func GetPublicacionPorId(c *gin.Context) {
	publicacion := &models.Publicacion{}
	id := c.Param("id")
	err := models.Db.Joins("Usuario").Joins("Etapa").First(publicacion, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.CrearRespuesta(errors.New("Publicacion no encontrada"), nil, c, http.StatusNotFound)
			return
		}
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener publicacion"), nil, c, http.StatusInternalServerError)
		return
	}
	for _, img := range publicacion.ImagenesPublicacion {
		if img.Imagen != "" {
			img.Imagen = utils.SERVIMG + img.Imagen
		} else {
			img.Imagen = utils.DefaultNoticia
		}
	}
	err = models.Db.Where("usuario_id = ? and publicacion_id = ?", c.GetInt("id_usuario"), publicacion.ID).First(&models.LecturaPublicacion{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al obtener publicaciones"), nil, c, http.StatusInternalServerError)
			return
		}
	}

	utils.CrearRespuesta(nil, publicacion, c, http.StatusOK)
}

type MediaPublicacion struct {
	Titulo   string `form:"titulo"`
	Imagenes string `form:"titulo"`
}

func CreatePublicacionMedia(c *gin.Context) {

	idParroquia := uint(c.GetInt("id_etapa"))
	publicacion := &models.Publicacion{}
	publicacion.ParroquiaID = idParroquia
	// }
	form, _ := c.MultipartForm()
	var err error
	files := form.File["archivos[]"]
	tx := models.Db.Begin()
	titulo := form.Value["titulo"]
	if titulo == nil {
		utils.CrearRespuesta(errors.New("Por favor ingrese titulo"), nil, c, http.StatusBadRequest)
		return
	}
	cuerpo := form.Value["cuerpo"]
	if cuerpo == nil {

		utils.CrearRespuesta(errors.New("Por favor ingrese decripcion"), nil, c, http.StatusBadRequest)
		return
	}
	err = tx.Create(publicacion).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear publicacion"), nil, c, http.StatusInternalServerError)
		return
	}
	valor := 0
	numVideos := 0
	numImagenes := 0
	for _, file := range files {
		isVideo := false
		fileArr := strings.Split(file.Filename, ".")
		extension := fileArr[len(fileArr)-1]
		if extension == "mp4" {
			isVideo = true
			numVideos++
			if numVideos == 2 {
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Solo puede subir un video por noticia"), nil, c, http.StatusBadRequest)
				return
			}
		} else if extension == "jpeg" || extension == "jpg" || extension == "png" {
			isVideo = false
			numImagenes++
			if numImagenes == 4 {
				tx.Rollback()
				utils.CrearRespuesta(errors.New("Solo se puede subir hasta 4 imagenes por noticia"), nil, c, http.StatusBadRequest)
				return
			}
		} else {
			tx.Rollback()
			utils.CrearRespuesta(errors.New("Archivos de formato ."+extension+" no son permitidos"), nil, c, http.StatusBadRequest)
			return
		}

		v := fmt.Sprintf("%d", valor)
		tiempo := fmt.Sprintf("%d", time.Now().Unix())
		nombre := "public/img/publicacions/" + tiempo + v + file.Filename
		err = c.SaveUploadedFile(file, nombre)
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error en formato de imagenes"), nil, c, http.StatusBadRequest)
			return
		}
		err = tx.Create(&models.PublicacionImagen{Imagen: nombre, PublicacionID: publicacion.ID, IsVideo: isVideo}).Error
		if err != nil {
			tx.Rollback()
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al guardar publicacion"), nil, c, http.StatusInternalServerError)
			return
		}
		valor++
		print(valor)
	}

	//Notificar emprendimientos
	residentes := []*models.Fiel{}
	err = models.Db.Where("Casa.etapa_id = ?", publicacion.ParroquiaID).Joins("Casa").Find(&residentes).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(err, "Error al crear publicacion", c, http.StatusInternalServerError)
		return
	}
	tokens := []string{}
	for _, res := range residentes {
		tokens = append(tokens, res.TokenNotificacion)
	}
	go notification.SendNotification("Nueva Publicacion", publicacion.Titulo, tokens, "1")
	tx.Commit()
	utils.CrearRespuesta(nil, "Publicacion creada exitosamente", c, http.StatusCreated)

}

func CreatePublicacion(c *gin.Context) {
	idParroquia := uint(c.GetInt("id_etapa"))
	publicacion := &models.Publicacion{}
	err := c.ShouldBindJSON(publicacion)
	publicacion.ParroquiaID = idParroquia
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	tx := models.Db.Begin()
	for _, imagen := range publicacion.ImagenesPublicacion {
		if imagen.Imagen != "" {

			imagen.Imagen, err = img.FromBase64ToImage(imagen.Imagen, "publicacions/"+time.Now().Format(time.RFC3339), false)
			if err != nil {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error en formato de imagenes"), nil, c, http.StatusOK)
				return
			}
		}
	}
	err = tx.Create(publicacion).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear publicacion"), nil, c, http.StatusInternalServerError)
		return
	}
	//Notificar emprendimientos
	residentes := []*models.Fiel{}
	err = models.Db.Where("Casa.etapa_id = ?", publicacion.ParroquiaID).Joins("Casa").Find(&residentes).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(err, "Error al crear publicacion", c, http.StatusInternalServerError)
		return
	}
	tokens := []string{}
	for _, res := range residentes {
		tokens = append(tokens, res.TokenNotificacion)
	}
	go notification.SendNotification("Nueva Publicacion", publicacion.Titulo, tokens, "1")
	tx.Commit()
	utils.CrearRespuesta(err, "Publicacion creada exitosamente", c, http.StatusCreated)

}

func UpdatePublicacion(c *gin.Context) {
	publicacion := &models.Publicacion{}

	err := c.ShouldBindJSON(publicacion)
	if err != nil {
		utils.CrearRespuesta(err, nil, c, http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	idPublicacion := uint(id)
	tx := models.Db.Begin()
	for _, ima := range publicacion.ImagenesPublicacion {
		ima.PublicacionID = idPublicacion
		ima.Imagen, err = img.FromBase64ToImage(ima.Imagen, "publicacions/"+time.Now().Format(time.RFC3339), false)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error"), nil, c, http.StatusBadRequest)
			return
		}
	}

	err = tx.Omit("publicacion", "usuario_id", "etapa_id").Where("id = ?", id).Updates(publicacion).Error
	if err != nil {
		tx.Rollback()
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al actualizar publicacion"), nil, c, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	utils.CrearRespuesta(err, "Publicacion actualizada correctamente", c, http.StatusOK)
}

func DeletePublicacion(c *gin.Context) {
	id := c.Param("id")
	err := models.Db.Delete(&models.Publicacion{}, id).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al eliminar publicacion"), nil, c, http.StatusInternalServerError)
		return
	}
	utils.CrearRespuesta(nil, "Publicacion eliminada exitosamente", c, http.StatusOK)
}

func MarcarPublicacionLeida(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	idPublicacion := uint(id)
	usuarioId := c.GetInt("id_usuario")
	err := models.Db.Create(&models.LecturaPublicacion{UsuarioID: uint(usuarioId), PublicacionID: idPublicacion}).Error
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al crear publicacion"), nil, c, http.StatusInternalServerError)
	}
	utils.CrearRespuesta(nil, "Lectura de noticia correcta", c, http.StatusOK)
}
