package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/drmendoz/iglesias-backend/models"
)

var JwtKey = []byte("my_secret_key")

type Login struct {
	Correo            string `json:"correo" binding:"required"`
	Contrasena        string `json:"contrasena" binding:"required"`
	TokenNotificacion string `json:"token_notificacion" `
}
type Claims struct {
	Rol         string
	Id          int
	TokenRandom string
	jwt.StandardClaims
}

func HashPassword(password string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(password))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

func GenerarToken(usuario *models.Usuario, rol string) string {

	expirationTime := time.Now().Add(500000 * time.Minute)
	claims := &Claims{
		Rol: rol,
		Id:  int(usuario.ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		TokenRandom: usuario.RandomNumToken,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(JwtKey)
	return tokenString
}
