package numgen

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
)

type NumGen struct {
}

func NewNumGen() *NumGen {
	return &NumGen{}
}

func (n *NumGen) GenCardNumber(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return gofakeit.CreditCard().Number, nil
}
