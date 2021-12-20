package main

import (
	"fmt"

	"github.com/drmendoz/iglesias-backend/routes"
	"github.com/drmendoz/iglesias-backend/utils"
)

func main() {
	port := ":8080"
	if utils.Viper.GetBool("PROD") {
		port = ":" + utils.Viper.GetString("APP_PORT")
	}
	err := routes.R.Run(port)
	if err != nil {
		str := fmt.Sprintf("Error al utilizar puerto %s verifique su dispponibilidad. Error: %v", port, err)
		utils.Log.Fatal(str)
	}

}
