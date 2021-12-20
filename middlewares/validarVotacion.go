package middlewares

import (
	"errors"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ValidarVoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		idVotacion := c.Param("id")
		user := c.GetInt("id_residente")
		res := &models.RespuestaVotacion{}
		err := models.Db.Where("residente_id = ? and OpcionVotacion.votacion_id = ? ", user, idVotacion).Joins("OpcionVotacion").First(res).Error
		utils.Log.Info(c.GetBool("is_principal"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if !c.GetBool("is_principal") {
					c.Set("usuario_voto", true)
				}
				c.Next()
			} else {
				_ = c.Error(err)
				utils.CrearRespuesta(errors.New("Error al obtener Votacion"), nil, c, http.StatusInternalServerError)
				c.Abort()
				return
			}
		} else {
			c.Set("usuario_voto", true)
			c.Next()
		}
	}
}
