package core

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

func ParseCurrency(currency string) Currency {
	switch currency {
	case CurrencyUSD.String():
		return CurrencyUSD
	case CurrencyEUR.String():
		return CurrencyEUR
	case CurrencyGBP.String():
		return CurrencyGBP
	case CurrencyRUB.String():
		return CurrencyRUB
	case CurrencyEGP.String():
		return CurrencyEGP
	default:
		return ""
	}
}

type CreateAccountParams struct {
	UserID   int64
	Name     string
	Currency Currency
}

type Account struct {
	ID       int64    `json:"id"`
	OwnerID  int64    `json:"user_id"`
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
