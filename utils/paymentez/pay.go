package paymentez

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	http "net/http"
	"time"

	"github.com/drmendoz/iglesias-backend/auth"
	"github.com/drmendoz/iglesias-backend/utils"
)

type TarjetasResponse struct {
	Tarjetas   []*Tarjeta `json:"cards"`
	Resultados int        `json:"result_size"`
}

type Tarjeta struct {
	Bin                  *string `json:"bin"`
	Status               *string `json:"status"`
	Token                *string `json:"token"`
	HolderName           *string `json:"holder_name"`
	AnoExpiracion        *string `json:"expiry_year"`
	MesExpiracion        *string `json:"expiry_month"`
	TransactionReference *string `json:"transaction_reference"`
	Tipo                 *string `json:"type"`
	Numero               *string `json:"number"`
}

type DeleteCardRequest struct {
	Tarjeta CardToken `json:"card"`
	User    Usuario   `json:"user"`
}

type CardToken struct {
	Token string `json:"token"`
}

type Usuario struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type DeleteRespuesta struct {
	Mensaje string `json:"message"`
}

func generateAuthToken() string {

	unixTimeStamp := time.Now().Unix()
	appCode := utils.PaymentezAppCode
	appKey := utils.PaymentezAppKey
	uniqTokenString := fmt.Sprintf("%s%d", appKey, unixTimeStamp)
	uniqTokenStringHash := auth.HashPassword(uniqTokenString)
	token := fmt.Sprintf("%s;%d;%s", appCode, unixTimeStamp, uniqTokenStringHash)
	token = b64.StdEncoding.EncodeToString([]byte(token))
	print(token)
	return token
}
func GetTarjetas(idFiel int64) (interface{}, error) {
	url := fmt.Sprintf("https://ccapi-stg.paymentez.com/v2/card/list?uid=%d", idFiel)
	client := &http.Client{}
	authToken := generateAuthToken()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Auth-Token", authToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 300 {
		var data map[string]interface{}
		_ = json.NewDecoder(res.Body).Decode(&data)
		return data, nil
	}
	tarjetas := &TarjetasResponse{}
	err = json.NewDecoder(res.Body).Decode(tarjetas)
	utils.Log.Warn("%v", tarjetas)
	return tarjetas, err

}

func DeleteTarjeta(idFiel int, tokenTarjeta string) (interface{}, error) {
	url := "https://ccapi-stg.paymentez.com/v2/card/delete/"
	client := &http.Client{}
	authToken := generateAuthToken()
	id := fmt.Sprintf("%d", idFiel)
	body := &DeleteCardRequest{Tarjeta: CardToken{Token: tokenTarjeta}, User: Usuario{Id: id}}
	jsonBody, _ := json.Marshal(body)
	bodyBytes := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("POST", url, bodyBytes)
	req.Header.Add("Auth-Token", authToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 300 {
		var data map[string]interface{}
		_ = json.NewDecoder(res.Body).Decode(&data)
		return data, nil
	}
	mensaje := &DeleteRespuesta{}
	err = json.NewDecoder(res.Body).Decode(mensaje)
	return mensaje, err
}

type CobroRequest struct {
	Usuario Usuario `json:"user"`
	Orden   Orden   `json:"order"`
	Tarjeta Tarjeta `json:"card"`
}

type Orden struct {
	Monto         float64 `json:"amount"`
	Descripcion   string  `json:"description"`
	Referencia    string  `json:"dev_reference"`
	Iva           float64 `json:"vat"`
	PorcentajeIva float64 `json:"tax_percentage"`
}

type CobroResponse struct {
	Transaccion Transaction `json:"transaction"`
	Tarjeta     Tarjeta     `json:"card"`
}

type Transaction struct {
	Status             string  `json:"status"`
	FechaPago          string  `json:"payment_date"`
	Monto              float64 `json:"amount"`
	CodigoAutorizacion string  `json:"authorization_code"`
	Installments       int     `json:"installments"`
	Referencia         int     `json:"dev_reference"`
	Mensaje            string  `json:"message"`
	CodigoCarrier      string  `json:"carrier_code"`
	Id                 string  `json:"id"`
	DetalleStatus      string  `json:"status_detail"`
}

func CobrarTarjeta(idUsuario string, correo string, totalCompra float64, descripcion string, idTransaccion string, iva float64, tokenTarjeta string) (*CobroResponse, error) {
	url := "https://ccapi-stg.paymentez.com/v2/transaction/debit"
	client := &http.Client{}
	authToken := generateAuthToken()
	body := &CobroRequest{Usuario: Usuario{Id: idUsuario, Email: correo}, Orden: Orden{Monto: totalCompra, Descripcion: descripcion, Referencia: idTransaccion, Iva: iva, PorcentajeIva: 0}, Tarjeta: Tarjeta{Token: &tokenTarjeta}}

	jsonBody, _ := json.Marshal(body)
	bodyBytes := bytes.NewBuffer(jsonBody)
	fmt.Printf("%v", body)
	req, _ := http.NewRequest("POST", url, bodyBytes)
	req.Header.Add("Auth-Token", authToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 300 {
		var data map[string]interface{}
		_ = json.NewDecoder(res.Body).Decode(&data)
		fmt.Printf("%v", data)
		return nil, errors.New("Error con el pago de tarjeta")
	}

	mensaje := &CobroResponse{}
	_ = json.NewDecoder(res.Body).Decode(mensaje)
	if mensaje.Transaccion.Status != "success" {
		return nil, errors.New("Error con el pago de tarjeta")
	}
	return mensaje, nil

}

func AnadirTarjeta(idFiel int, tokenTarjeta string) (interface{}, error) {
	url := "https://ccapi-stg.paymentez.com/v2/card/delete/"
	client := &http.Client{}
	authToken := generateAuthToken()
	id := fmt.Sprintf("%d", idFiel)
	body := &DeleteCardRequest{Tarjeta: CardToken{Token: tokenTarjeta}, User: Usuario{Id: id}}
	jsonBody, _ := json.Marshal(body)
	bodyBytes := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("POST", url, bodyBytes)
	req.Header.Add("Auth-Token", authToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 300 {
		var data map[string]interface{}

		fmt.Printf("%v", data)
		_ = json.NewDecoder(res.Body).Decode(&data)
		return data, errors.New("Error al debitar tarjeta")
	}
	mensaje := &DeleteRespuesta{}
	err = json.NewDecoder(res.Body).Decode(mensaje)
	return mensaje, err
}

type DevolverBody struct {
	Transaction *Transaction `json:"transaction"`
}

type Status struct {
	Status  string `json:"status"`
	Detalle string `json:"detail"`
}

func DevolverPago(idTransacion uint) (*Status, error) {
	url := "https://ccapi-stg.paymentez.com/v2/transaction/refund/"
	client := &http.Client{}
	authToken := generateAuthToken()
	id := fmt.Sprintf("%d", idTransacion)
	body := &DevolverBody{Transaction: &Transaction{Id: id}}
	jsonBody, _ := json.Marshal(body)
	bodyBytes := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("POST", url, bodyBytes)
	req.Header.Add("Auth-Token", authToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	mensaje := &Status{}
	err = json.NewDecoder(res.Body).Decode(mensaje)
	return mensaje, nil
}
