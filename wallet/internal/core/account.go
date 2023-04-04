package core

import (
	"errors"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyRUB Currency = "RUB"
	CurrencyEGP Currency = "EGP"
)

func (c Currency) String() string {
	return string(c)
}

var (
	ErrInvalidCurrency = errors.New("invalid currency")
)

func ParseCurrency(currency string) (Currency, error) {
	switch currency {
	case "USD":
		return CurrencyUSD, nil
	case "EUR":
		return CurrencyEUR, nil
	case "GBP":
		return CurrencyGBP, nil
	case "RUB":
		return CurrencyRUB, nil
	case "EGP":
		return CurrencyEGP, nil
	default:
		return "", ErrInvalidCurrency
	}
}

type CreateAccountParams struct {
	UserID   int64
	Name     string
	Currency Currency
}

type Account struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Balance  float64  `json:"balance"`
	Currency Currency `json:"currency"`
}

type CreateCardParams struct {
	AccountID int64
	Number    string
}

type Card struct {
	AccountID int64  `json:"account_id"`
	Number    string `json:"number"`
}
