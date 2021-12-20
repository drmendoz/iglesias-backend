package models

import "gorm.io/gorm"

type PagoMarket struct {
	gorm.Model
	ValorPagado  float64        `json:"valor_pagado"`
	FechaDePago  float64        `json:"fecha_pago"`
	ItemMarketID uint           `json:"id_item_market"`
	ItemMarket   Emprendimiento `json:"item_market"`
}
