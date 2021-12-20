package img

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/drmendoz/iglesias-backend/utils"
)

func FromBase64ToImage(base64img string, nombre string, recortar bool) (string, error) {
	if !strings.Contains(base64img, ",") {
		return "", errors.New("Format: Error en formato de imagen")
	}
	listBase64 := strings.Split(base64img, ",")
	formato := listBase64[0]
	img := listBase64[1]
	unbased, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(unbased)
	var im image.Image
	extension := ""
	if formato == "data:image/png;base64" {
		im, err = png.Decode(r)
		extension = ".png"

	} else if formato == "data:image/jpg;base64" || formato == "data:image/jpeg;base64" {
		im, err = jpeg.Decode(r)
		extension = ".jpeg"
	} else {
		return "", errors.New("Decode: Error en formato de imagen")
	}

	if err != nil {
		utils.Log.Warn(err)
		return "", errors.New("Error al decodificar iamgen")
	}
	// if recortar {
	// 	im, err = cutter.Crop(im, cutter.Config{
	// 		Width:  300,
	// 		Height: 300,
	// 		Mode:   cutter.Centered,
	// 	})
	// 	if err != nil {
	// 		return "", err
	// 	}
	// }
	path := "public/img/" + nombre + extension

	path = strings.ReplaceAll(path, " ", "-")
	path = strings.ReplaceAll(path, ":", "-")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	if extension == ".png" {
		err = png.Encode(f, im)
	} else {
		err = jpeg.Encode(f, im, nil)
	}

	if err != nil {
		return "", errors.New("Error al guardar imagen")
	}
	return path, err
}

func IsBase64(str string) bool {
	res := strings.HasPrefix(str, "data:image/")
	return res
}
