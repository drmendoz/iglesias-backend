package utils

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type LocalFormatter struct {
	log.Formatter
}

var Log *log.Logger = log.New()

func init() {
	file, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		println("No existe carpeta logs")
		os.Exit(1)
	}
	Log.SetFormatter(LocalFormatter{&log.TextFormatter{}})
	Log.SetOutput(io.MultiWriter(os.Stdout, file))
}

func (u LocalFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.Local()
	return u.Formatter.Format(e)
}
