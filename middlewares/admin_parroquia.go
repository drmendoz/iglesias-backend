package middlewares

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func ParsingTokenAdminParroquia() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetInt("id_usuario")
		adm := &models.AdminParroquia{}
		res := models.Db.Where("Usuario.id = ? ", user).Joins("Usuario").Joins("Parroquia").First(adm)
		if res.Error != nil {
			utils.CrearRespuesta(errors.New("Error de autorizacion"), nil, c, http.StatusInternalServerError)
			c.Abort()
			return
		}
		if adm.Parroquia == nil {
			utils.CrearRespuesta(errors.New("Su parroquia ya no existe"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("id_parroquia", int(adm.ParroquiaID))
		c.Next()
	}
}
