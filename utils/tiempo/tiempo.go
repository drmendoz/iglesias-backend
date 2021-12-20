package tiempo

import (
	"time"

	"github.com/drmendoz/iglesias-backend/utils"
)

var Local *time.Location

func init() {
	loc, err := time.LoadLocation("America/Guayaquil")
	if err != nil {
		utils.Log.Fatal("Error al obtener localizacion")

	}
	Local = loc
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
