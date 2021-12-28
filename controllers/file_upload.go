package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/drmendoz/iglesias-backend/utils/img"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadImagePerfil(usuario *models.Usuario, tx *gorm.DB) (*models.Usuario, error) {
	if usuario.Imagen != "" {
		usuarioID := fmt.Sprintf("%d", usuario.ID)

		img, err := img.FromBase64ToImage(usuario.Imagen, "usuarios/"+time.Now().Format(time.RFC3339)+usuarioID, false)
		if err != nil {
			return usuario, err
		}
		usuario.Imagen = img
		err = tx.Model(&models.Usuario{}).Where("id = ?", usuario.ID).Update("imagen", img).Error
		return usuario, err
	}
	return usuario, nil
}

func SubirArchivos(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		_ = c.Error(err)
		utils.CrearRespuesta(errors.New("Error al obtener archivos"), nil, c, http.StatusBadRequest)
		return
	}
	files := form.File["archivos[]"]
	num := 0

	archivos := []string{}
	for _, file := range files {

		v := fmt.Sprintf("%d", num)
		tiempo := fmt.Sprintf("%d", time.Now().Unix())
		fileArr := strings.Split(file.Filename, ".")
		extension := fileArr[len(fileArr)-1]
		nombre := "public/img/emprendimiento/" + tiempo + v + "." + extension
		err = c.SaveUploadedFile(file, nombre)
		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al subir archivos adjuntos"), nil, c, http.StatusBadRequest)
			return
		}

		if err != nil {
			_ = c.Error(err)
			utils.CrearRespuesta(errors.New("Error al guardar publicacion"), nil, c, http.StatusInternalServerError)
			return
		}
		archivos = append(archivos, nombre)
		num++
	}
	utils.CrearRespuesta(nil, archivos, c, http.StatusOK)

}
