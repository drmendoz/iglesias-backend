package controllers

import (
	"fmt"
	"time"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils/img"
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
