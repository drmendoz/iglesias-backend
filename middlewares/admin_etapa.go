package middlewares

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func ParsingTokenAdminEtapa() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetInt("id_usuario")
		adm := &models.AdminEtapa{}
		res := models.Db.Where("Usuario.id = ? ", user).Joins("Usuario").Joins("Etapa").First(adm)
		if res.Error != nil {
			utils.CrearRespuesta(errors.New("Error de autorizacion"), nil, c, http.StatusInternalServerError)
			c.Abort()
			return
		}
		if adm.Etapa == nil {
			utils.CrearRespuesta(errors.New("Su etapa ya no existe"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		print(adm.EtapaID)
		c.Set("id_etapa", int(adm.EtapaID))
		c.Next()
	}
}
