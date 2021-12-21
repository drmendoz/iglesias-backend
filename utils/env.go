package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

var SERVIMG string
var DefaultUrb string
var DefaultUser string
var DefaultNoticia string
var DefaultParroquia string
var DefaultAdministrativo string
var DefaultVisita string
var DefaultGaleria string
var DefaultExpreso string
var DefaultVotacion string
var DefaultMensaje string
var DefaultPublicidad string
var DefaultAreaSocial string
var DefaultCasa string
var DefaultCam string
var DefaultCategoria string
var Colores = map[int]string{1: "#7875EA", 2: "#6C69D1", 3: "#5855AB", 4: "#37366B"}
var CamServer string
var NumeroMaximoEmprendimiento int
var RutaTutorial string
var PaymentezAppCode string
var PaymentezAppKey string
var LogoPractical string

func init() {
	Viper = viper.New()
	Viper.SetConfigFile(".env")
	err := Viper.ReadInConfig()
	if err != nil {
		fmt.Printf("No existe el archivo .env. Error: %v ", err)
		os.Exit(1)
	}
	SERVIMG = Viper.GetString("SERVER_DEV_IMG")
	if Viper.GetBool("PROD") {
		SERVIMG = Viper.GetString("SERVER_IMG")
	}
	NumeroMaximoEmprendimiento = 6
	DefaultNoticia = Viper.GetString("DEFAULT_NOTICIA")
	DefaultUser = Viper.GetString("DEFAULT_USER")
	DefaultUrb = Viper.GetString("DEFAULT_URB")
	DefaultParroquia = Viper.GetString("DEFAULT_PARROQUIA")
	DefaultAdministrativo = Viper.GetString("DEFAULT_ADMINISTRATIVO")
	DefaultVisita = Viper.GetString("DEFAULT_VISITA")
	DefaultGaleria = Viper.GetString("DEFAULT_GALERIA")
	DefaultPublicidad = Viper.GetString("DEFAULT_PUBLICIDAD")
	DefaultExpreso = Viper.GetString("DEFAULT_EXPRESO")
	DefaultVotacion = Viper.GetString("DEFAULT_VOTACION")
	DefaultMensaje = Viper.GetString("DEFAULT_MENSAJE")
	DefaultAreaSocial = Viper.GetString("DEFAULT_AREA")
	DefaultCasa = Viper.GetString("DEFAULT_CASA")
	CamServer = Viper.GetString("CAM_SERVER")
	DefaultCam = Viper.GetString("DEFAULT_CAM")
	DefaultCategoria = Viper.GetString("DEFAULT_CAT")
	RutaTutorial = Viper.GetString("RUTA_TUTORIAL")
	PaymentezAppCode = Viper.GetString("PAYMENTEZ_APP_CODE")
	PaymentezAppKey = Viper.GetString("PAYMENTEZ_APP_KEY")
	LogoPractical = Viper.GetString("LOGO_PRACTICAL")

}
