package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

func ParsingTokenResidente() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetInt("id_usuario")
		adm := &models.Residente{}
		err := models.Db.Where("Usuario.id = ? ", user).Joins("Usuario").Preload("Casa").Preload("Casa.Etapa").Preload("Casa.Etapa.Urbanizacion").First(adm).Error
		if err != nil {
			utils.CrearRespuesta(errors.New("Error de autorizacion"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		if adm.Casa == nil || adm.Casa.Etapa == nil || adm.Casa.Etapa.Urbanizacion == nil {
			utils.CrearRespuesta(errors.New("Usuario no permitido"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenRandom := c.GetString("token_random")
		if tokenRandom == "" || tokenRandom != adm.Usuario.RandomNumToken {
			utils.CrearRespuesta(errors.New("Error de autorizacion"), nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		fmt.Printf("Token Random: %s\n", tokenRandom)
		fmt.Printf("%d", adm.Casa.EtapaID)
		c.Set("id_residente", int(adm.ID))
		c.Set("id_casa", int(adm.Casa.ID))
		c.Set("id_etapa", int(adm.Casa.EtapaID))
		c.Set("id_urbanizacion", int(adm.Casa.Etapa.UrbanizacionID))
		c.Set("is_principal", adm.IsPrincipal)
		fmt.Printf("%d", int(adm.Casa.Etapa.UrbanizacionID))
		c.Next()
	}
}
