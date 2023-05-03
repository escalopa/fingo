package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCurrency(t *testing.T) {

	tests := []struct {
		name     string
		currency string
		want     Currency
	}{
		{
			name:     "USD",
			currency: "USD",
			want:     CurrencyUSD,
		},
		{
			name:     "EUR",
			currency: "EUR",
			want:     CurrencyEUR,
		},
		{
			name:     "GBP",
			currency: "GBP",
			want:     CurrencyGBP,
		},
		{
			name:     "RUB",
			currency: "RUB",
			want:     CurrencyRUB,
		},
		{
			name:     "EGP",
			currency: "EGP",
			want:     CurrencyEGP,
		},
		{
			name:     "empty",
			currency: "",
			want:     "",
		},
		{
			name:     "unknown",
			currency: "unknown",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ParseCurrency(tt.currency))
		})
	}
}
