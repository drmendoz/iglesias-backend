package routes

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func setLogging() {
	f, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print(err)
	}
	fError, err := os.OpenFile("logs/gin_error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print(err)
	}

	gin.ForceConsoleColor()

	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stdout, fError)

}
