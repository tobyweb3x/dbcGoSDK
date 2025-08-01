package maths

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/types"
	"errors"
	"math/big"
)

func MulDiv(x, y, denominator *big.Int, rounding types.Rounding) (*big.Int, error) {
	if denominator.Sign() == 0 {
		return nil, errors.New("MulDiv: division by zero")
	}

	prod := new(big.Int).Mul(x, y)
	if denominator.Cmp(big.NewInt(1)) == 0 || x.Sign() == 0 || y.Sign() == 0 {
		return prod, nil
	}

	if rounding == types.RoundingUp {
		// Calculate ceiling division: (prod + denominator - 1) / denominator
		return new(big.Int).Quo(
			new(big.Int).Add(
				prod,
				new(big.Int).Sub(denominator, big.NewInt(1)),
			),
			denominator,
		), nil
	}

	return new(big.Int).Quo(prod, denominator), nil
}

func Q64(n float64) *big.Int {
	f := new(big.Float).Mul(
		new(big.Float).SetFloat64(n),
		new(big.Float).SetFloat64(1<<64),
	)
	i := new(big.Int)
	f.Int(i) // truncates
	return i
}

// PowQ64 is a custom power function for [big.Int] with scaling.
func PowQ64(base, exponent *big.Int, scaling bool) *big.Int {
	// result := new(big.Int).Set(oneQ64)

	// special cases
	if exponent.Sign() == 0 {
		return constants.OneQ64
	}
	if base.Sign() == 0 {
		return big.NewInt(0)
	}
	if base.Cmp(constants.OneQ64) == 0 {
		return constants.OneQ64
	}

	isNegative, absExponent := exponent.Sign() < 0, new(big.Int).Abs(exponent)

	result, currentBase, exp := constants.OneQ64, base, absExponent
	for exp.Sign() != 0 {
		if one := big.NewInt(1); new(big.Int).And(exp, one).Cmp(one) == 0 {
			result = new(big.Int).Quo(
				new(big.Int).Mul(result, currentBase),
				constants.OneQ64,
			)
		}
		currentBase = new(big.Int).Quo(
			new(big.Int).Mul(currentBase, currentBase),
			constants.OneQ64,
		)
		exp = new(big.Int).Rsh(exp, 1)
	}

	// handle negative exponent
	if isNegative {
		result = new(big.Int).Div(
			new(big.Int).Mul(constants.OneQ64, constants.OneQ64),
			result,
		)
	}

	if !scaling {
		return new(big.Int).Quo(result, constants.OneQ64)
	}

	return result
}
