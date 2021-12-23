package middlewares

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func ParsingTokenFiel() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetInt("id_usuario")
		adm := &models.Fiel{}
		err := models.Db.Where("Usuario.id = ? ", user).Joins("Usuario").Preload("Parroquia").First(adm).Error
		if err != nil {
			utils.CrearRespuesta(errors.New("Error de autorizacion"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		tokenRandom := c.GetString("token_random")
		if tokenRandom == "" || tokenRandom != adm.Usuario.RandomNumToken {
			utils.CrearRespuesta(errors.New("Error de autorizacion"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("id_fiel", int(adm.ID))
		c.Set("id_parroquia", int(*adm.ParroquiaID))
		c.Next()
	}
}
