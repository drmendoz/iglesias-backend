package utils

import (
	"encoding/base64"
	"os"
)

func SubirPdf(nombre string, contenido string) error {
	dec, err := base64.StdEncoding.DecodeString(contenido)
	if err != nil {
		return err
	}

	f, err := os.Create("public/pdf/" + nombre)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return err
}
