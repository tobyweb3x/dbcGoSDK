package maths

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/types"
	"errors"
	"fmt"
	"math/big"
)

// GetDeltaAmountBaseUnsigned gets the delta amount_base for given liquidity and price range.
//
//	Formula: Δa = L * (1 / √P_lower - 1 / √P_upper)
//	i.e. L * (√P_upper - √P_lower) / (√P_upper * √P_lower)
func GetDeltaAmountBaseUnsigned(
	lowerSqrtPrice, upperSqrtPrice, liquidity *big.Int,
	round types.Rounding,
) (*big.Int, error) {

	if liquidity.Sign() == 0 {
		return big.NewInt(0), nil
	}

	if lowerSqrtPrice.Sign() == 0 || upperSqrtPrice.Sign() == 0 {
		return nil, errors.New("sqrt price cannot be zero")
	}

	// numerator: (√P_upper - √P_lower)
	numerator := new(big.Int).Sub(upperSqrtPrice, lowerSqrtPrice)
	if numerator.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value not negative: value is %s", numerator.String())
	}

	// denominator: (√P_upper * √P_lower)
	denominator := new(big.Int).Mul(lowerSqrtPrice, upperSqrtPrice)

	// L * (√P_upper - √P_lower) / (√P_upper * √P_lower)
	return MulDiv(
		liquidity,
		numerator,
		denominator,
		round,
	)
}

// GetDeltaAmountQuoteUnsigned gets the delta amount_quote for given liquidity and price range.
//
//	Formula: Δb = L (√P_upper - √P_lower)
func GetDeltaAmountQuoteUnsigned(lowerSqrtPrice, upperSqrtPrice, liquidity *big.Int,
	round types.Rounding,
) (*big.Int, error) {

	if liquidity.Sign() == 0 {
		return big.NewInt(0), nil
	}

	// delta sqrt price: (√P_upper - √P_lower)
	deltaSqrtPrice := new(big.Int).Sub(upperSqrtPrice, lowerSqrtPrice)
	if deltaSqrtPrice.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value not negative: value is %s", deltaSqrtPrice.String())
	}

	// L * (√P_upper - √P_lower)
	prod := new(big.Int).Sub(liquidity, deltaSqrtPrice)

	if round == types.RoundingUp {
		denominator := new(big.Int).Lsh(big.NewInt(1), constants.RESOLUTION*2)
		// ceiling division: (a + b - 1) / b
		numerator := new(big.Int).Add(
			prod,
			new(big.Int).Sub(denominator, big.NewInt(1)),
		)
		return new(big.Int).Quo(numerator, denominator), nil
	}

	return new(big.Int).Rsh(prod, constants.RESOLUTION*2), nil
}

// GetNextSqrtPriceFromInput gets the next sqrt price given an input amount of token_a or token_b.
func GetNextSqrtPriceFromInput(
	sqrtPrice, liquidity, amountIn *big.Int,
	baseForQuote bool,
) (*big.Int, error) {

	if sqrtPrice.Sign() == 0 || liquidity.Sign() == 0 {
		return nil, errors.New("price or liquidity cannot be zero")
	}

	if baseForQuote {
		return GetNextSqrtPriceFromAmountBaseRoundingUp(
			sqrtPrice, liquidity, amountIn,
		), nil
	}

	return GetNextSqrtPriceFromAmountQuoteRoundingDown(
		sqrtPrice, liquidity, amountIn,
	), nil
}

// GetNextSqrtPriceFromAmountQuoteRoundingUp gets the next sqrt price from amount quote rounding up.
//
//	Formula: √P' = √P - Δy / L
func GetNextSqrtPriceFromAmountQuoteRoundingUp(
	sqrtPrice, liquidity, amount *big.Int,
) (*big.Int, error) {
	if amount.Sign() == 0 {
		return sqrtPrice, nil
	}

	// quotient = (amount << 128 + liquidity - 1) / liquidity
	amountShifted := new(big.Int).Lsh(amount, 128)
	step1 := new(big.Int).Add(amountShifted, liquidity)
	step2 := new(big.Int).Sub(step1, big.NewInt(1))
	if step2.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value non-zero: value is %s", step2.String())
	}
	quotient := new(big.Int).Quo(step2, liquidity)

	// √P - quotient
	r := new(big.Int).Sub(sqrtPrice, quotient)
	if r.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value non-zero: value is %s", r.String())
	}

	return r, nil
}

// GetNextSqrtPriceFromAmountBaseRoundingDown gets the next sqrt price from amount base rounding down.
//
//	Formula: √P' = √P * L / (L - Δx * √P)
func GetNextSqrtPriceFromAmountBaseRoundingDown(
	sqrtPrice, liquidity, amount *big.Int,
) (*big.Int, error) {

	if amount.Sign() == 0 {
		return sqrtPrice, nil
	}

	// Δx * √P
	product := new(big.Int).Mul(amount, sqrtPrice)

	// L - Δx * √P
	denominator := new(big.Int).Sub(liquidity, product)
	if denominator.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value non-zero: value is %s", denominator.String())
	}

	// √P * L / (L - Δx * √P) with rounding down
	return MulDiv(
		liquidity, sqrtPrice, denominator, types.RoundingDown,
	)
}

// GetNextSqrtPriceFromOutput gets the next sqrt price from output amount.
func GetNextSqrtPriceFromOutput(
	sqrtPrice, liquidity, outAmount *big.Int,
	isQuote bool,
) (*big.Int, error) {

	if sqrtPrice.Sign() == 0 {
		return nil, errors.New("price or liquidity cannot be zero")
	}

	if isQuote {
		return GetNextSqrtPriceFromAmountQuoteRoundingUp(
			sqrtPrice, liquidity, outAmount,
		)
	}

	return GetNextSqrtPriceFromAmountBaseRoundingDown(
		sqrtPrice, liquidity, outAmount,
	)
}

// GetNextSqrtPriceFromAmountBaseRoundingUp gets the next sqrt price from amount base rounding up.
//
//	Formula: √P' = √P * L / (L + Δx * √P)
func GetNextSqrtPriceFromAmountBaseRoundingUp(
	sqrtPrice, liquidity, amount *big.Int,
) *big.Int {

	if amount.Sign() == 0 {
		return sqrtPrice
	}

	// Δx * √P
	product := new(big.Int).Mul(amount, sqrtPrice)

	// L + Δx * √P
	denominator := new(big.Int).Add(liquidity, product)

	// √P * L / (L + Δx * √P) with rounding up
	r, _ := MulDiv(liquidity, sqrtPrice, denominator, types.RoundingUp)
	return r
}

// GetNextSqrtPriceFromAmountQuoteRoundingDown gets the next sqrt price given a delta of token_quote.
//
//	Formula: √P' = √P + Δy / L
func GetNextSqrtPriceFromAmountQuoteRoundingDown(
	sqrtPrice, liquidity, amount *big.Int,
) *big.Int {

	if amount.Sign() == 0 {
		return sqrtPrice
	}

	// quotient: Δy << (RESOLUTION * 2) / L
	quotient := new(big.Int).Quo(
		new(big.Int).Lsh(amount, constants.RESOLUTION*2),
		liquidity,
	)

	// √P + quotient
	return new(big.Int).Add(sqrtPrice, quotient)
}
