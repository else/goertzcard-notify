package client

import "github.com/shopspring/decimal"

type GoertzClient struct {
	username string
	password string
	sessId   string
	errChan  chan (error)
}

type GoertzCard struct {
	client      *GoertzClient
	kn          string
	Ean         string
	Description string
	Balance     decimal.Decimal
}
