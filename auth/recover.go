package auth

import (
	"crypto/rand"
	"io"
)

type NuevaContrasena struct {
	Contrasena     string `json:"nueva_contrasena"`
	Correo         string `json:"correo"`
	CodigoTemporal string `json:"codigo_temporal"`
	Imagen         string `json:"imagen"`
}

func GenerarCodigoTemporal(longitud int) (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, longitud)
	n, err := io.ReadAtLeast(rand.Reader, b, longitud)
	if n != longitud {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), err

}
