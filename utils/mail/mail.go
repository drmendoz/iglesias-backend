package mail

import (
	"net/smtp"

	"github.com/drmendoz/iglesias-backend/models"
	"github.com/drmendoz/iglesias-backend/utils"
)

var usuario string
var contrasena string

var auth smtp.Auth

func init() {
	usuario = utils.Viper.GetString("MAIL_USER")
	contrasena = utils.Viper.GetString("MAIL_PASSWORD")
	auth = smtp.PlainAuth("", usuario, contrasena, "smtp.gmail.com")
}

func EnviarCorreoRecoverContraseña(usuario models.Usuario) error {
	r := NewRequest([]string{usuario.Correo}, "Practical App: Creación de cuenta", "Creación de cuenta!")
	err := r.ParseTemplate("views/recover.html", usuario)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}

func EnviarCorreoRecover(usuario models.Usuario) error {
	r := NewRequest([]string{usuario.Correo}, "Practical App: Recuperacion de cuenta", "Recuperacion de cuenta!")
	err := r.ParseTemplate("views/recover.html", usuario)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}

func EnviarCambioContrasena(usuario models.Fiel) error {
	r := NewRequest([]string{usuario.Usuario.Correo}, "Practical App: Cambio de contrasena obligatorio", "Cambio de contrasena!")
	err := r.ParseTemplate("views/cambio.html", usuario)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}

func EnviarCambioContrasenaMaster(usuario models.AdminMaster) error {
	r := NewRequest([]string{usuario.Usuario.Correo}, "Practical App: Cambio de contrasena obligatorio", "Cambio de contrasena!")
	err := r.ParseTemplate("views/cambio.html", usuario)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}

func EnviarCambioContrasenaEtapa(usuario models.AdminParroquia) error {
	r := NewRequest([]string{usuario.Usuario.Correo}, "Practical App: Cambio de contrasena obligatorio", "Cambio de contrasena!")
	err := r.ParseTemplate("views/cambio.html", usuario)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}

func EnviarCambioContrasenaGarita(usuario models.AdminGarita) error {
	r := NewRequest([]string{usuario.Usuario.Correo}, "Practical App: Cambio de contrasena obligatorio", "Cambio de contrasena!")
	err := r.ParseTemplate("views/cambio.html", usuario)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}

type Contacto struct {
	Nombre      string `json:"nombre"`
	Correo      string `json:"correo"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
}

func EnviarCorreoContactoWeb(contacto Contacto) error {
	correoPractical := "info@practical.com.ec"
	r := NewRequest([]string{correoPractical}, "Practical: Solicitud de contacto", "Solicitud de contacto")
	err := r.ParseTemplate("views/contacto.html", contacto)
	if err != nil {
		return err
	}
	_, err = r.SendEmail()
	if err != nil {
		return err
	}
	return err
}
