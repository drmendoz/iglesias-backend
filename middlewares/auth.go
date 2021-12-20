package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/drmendoz/iglesias-backend/auth"
	"github.com/drmendoz/iglesias-backend/utils"
	"github.com/gin-gonic/gin"
)

//Credentials to login
type Credentials struct {
	Contrasena        string `json:"contrasena"`
	Correo            string `json:"correo"`
	TokenNotificacion string `json:"token_notificacion,omitempty"`
}

//Claims lo que se guarda en los tokens

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			err := errors.New("No posee token de autorizacion")
			_ = c.Error(err)
			utils.CrearRespuesta(err, nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		fmt.Printf("Token: %s\n", token)

		claims := &auth.Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if claims.Rol != c.GetString("rol") {
				err := errors.New("Rol no permitido")
				_ = c.Error(err)
				return nil, err
			}
			return auth.JwtKey, nil
		})
		if err != nil || err == jwt.ErrSignatureInvalid {
			err := c.Error(errors.New("Token Invalido"))
			_ = c.Error(err)
			utils.CrearRespuesta(err, nil, c, http.StatusUnauthorized)
			c.Abort()
			return

		}
		if !tkn.Valid {
			err := c.Error(errors.New("Token Expirado"))
			_ = c.Error(err)
			utils.CrearRespuesta(err, nil, c, http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Set("id_usuario", claims.Id)
		if claims.Rol == "fiel" {

			c.Set("token_random", claims.TokenRandom)
		}

		c.Next()

	}
}
