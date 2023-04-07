package numgen

import (
	"context"

	"github.com/brianvoe/gofakeit/v6"
)

type NumGen struct {
	l int // length of the generated number
}

func NewNumGen(l int) *NumGen {
	return &NumGen{l: l}
}

func (n *NumGen) GenCardNumber(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	var number []byte
	for i := 0; i < n.l; i++ {
		number = append(number, byte(gofakeit.Number(0, 9)+'0'))
	}
	return string(number), nil
}
