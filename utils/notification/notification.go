package notification

import (
	"github.com/drmendoz/iglesias-backend/utils"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func SendNotification(titulo string, cuerpo string, tokens []string, navigation string) {
	pushTokens := []expo.ExponentPushToken{}
	for _, token := range tokens {
		if token != "" {
			pushToken, _ := expo.NewExponentPushToken(token)
			pushTokens = append(pushTokens, pushToken)

		}
	}
	client := expo.NewPushClient(nil)
	_, err := client.Publish(
		&expo.PushMessage{
			To: pushTokens, Body: cuerpo,
			Data:     map[string]string{"tipo": navigation},
			Sound:    "default",
			Title:    titulo,
			Priority: expo.HighPriority},
	)
	if err != nil {
		utils.Log.Warn("Error al enviar notificaciones" + err.Error())
	}
}
