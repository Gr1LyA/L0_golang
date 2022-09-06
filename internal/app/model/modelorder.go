package model

import "time"

type ModelOrder struct {
	OrderUID    string `json:"order_uid" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Entry       string `json:"entry" validate:"required"`
	Delivery    struct {
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required"`
		Zip     string `json:"zip" validate:"required"`
		City    string `json:"city" validate:"required"`
		Address string `json:"address" validate:"required"`
		Region  string `json:"region" validate:"required"`
		Email   string `json:"email" validate:"required,email"`
	} `json:"delivery" validate:"required"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency" validate:"required"`
		Provider     string `json:"provider" validate:"required"`
		Amount       int    `json:"amount"`
		PaymentDt    int    `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items []struct {
		ChrtID      int    `json:"chrt_id" validate:"required"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		Rid         string `json:"rid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Sale        int    `json:"sale"`
		Size        string `json:"size" validate:"required"`
		TotalPrice  int    `json:"total_price"`
		NmID        int    `json:"nm_id"`
		Brand       string `json:"brand" validate:"required"`
		Status      int    `json:"status"`
	} `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
}