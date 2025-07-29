package maths

import (
	"dbcGoSDK/types"
	"errors"
	"math/big"
)

func MulDiv(x, y, denominator *big.Int, rounding types.Rounding) (*big.Int, error) {
	if denominator.Sign() == 0 {
		return nil, errors.New("MulDiv: division by zero")
	}

	if denominator.Cmp(big.NewInt(1)) == 0 || x.Sign() == 0 || y.Sign() == 0 {
		return new(big.Int).Mul(x, y), nil
	}

	prod := new(big.Int).Mul(x, y)

	if rounding == types.RoundingUp {
		// Calculate ceiling division: (prod + denominator - 1) / denominator
		numerator := new(big.Int).Add(
			prod,
			new(big.Int).Sub(denominator, big.NewInt(1)),
		)
		return numerator.Div(numerator, denominator), nil
	}

	return prod.Div(prod, denominator), nil
}
