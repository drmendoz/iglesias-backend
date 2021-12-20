package utils

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func EnviarSMS(telefono string, pin string, nombres string, apellidos string, mz string, villa string, etapa string, estado string) {
	println(telefono)
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: Viper.GetString("TWILIO_ACCOUNT_SID"),
		Password: Viper.GetString("TWILIO_AUTH_TOKEN"),
	})
	var mensaje string
	if estado != "RECHAZADA" {
		mensaje = fmt.Sprintf("PRACTICAL\n %s %s has recibido una autorizacion para ingresar a la Mz. %s Villa %s de %s, por favor al ser atendido por el personal de seguridad dictale el siguiente PIN: %s.\n Recuerda que el mismo pin te servira para salir.", nombres, apellidos, mz, villa, etapa, pin)
	} else {
		mensaje = fmt.Sprintf("PRACTICAL\n %s %s tu autorizacion de ingreso a la Mz. %s Villa %s de %s fue ANULADA. Por favor comunicate con la persona que te genero la autorizacion para que te genere una nueva.", nombres, apellidos, mz, villa, etapa)
	}
	params := &openapi.CreateMessageParams{}
	params.SetTo(telefono)
	params.SetFrom(Viper.GetString("TWILIO_PHONE_NUMBER"))
	params.SetBody(mensaje)

	_, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("SMS sent successfully!")
	}
}
