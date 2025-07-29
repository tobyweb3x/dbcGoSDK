package maths

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/types"
	"errors"
	"math/big"
)

// GetMaxIndex calculates the max index for rate limiter.
func GetMaxIndex(cliffFeeNumerator, feeIncrementBps *big.Int) *big.Int {
	deltaNumerator := new(big.Int).Sub(
		new(big.Int).SetUint64(constants.MaxFeeNumerator), cliffFeeNumerator)

	feeIncrementNumerator, _ := MulDiv(
		feeIncrementBps,
		new(big.Int).SetInt64(constants.FeeDenominator),
		new(big.Int).SetInt64(constants.BasisPointMax),
		types.RoundingDown,
	)

	return deltaNumerator.Div(deltaNumerator, feeIncrementNumerator)
}

// GetFeeNumeratorOnRateLimiter calculates the fee numerator on rate limiter.
func GetFeeNumeratorOnRateLimiter(
	cliffFeeNumerator, referenceAmount, feeIncrementBps, inputAmount *big.Int,
) (*big.Int, error) {

	if inputAmount.Cmp(referenceAmount) <= 0 {
		return cliffFeeNumerator, nil
	}

	if inputAmount.Sign() == 0 {
		return nil, errors.New("input amount cannot be zero")
	}

	diff := new(big.Int).Sub(inputAmount, referenceAmount)
	a, b := new(big.Int).QuoRem(diff, referenceAmount, new(big.Int))
	maxIndex := GetMaxIndex(cliffFeeNumerator, feeIncrementBps)
	i, _ := MulDiv(
		feeIncrementBps,
		new(big.Int).SetInt64(constants.FeeDenominator),
		new(big.Int).SetInt64(constants.BasisPointMax),
		types.RoundingDown,
	)

	one, two, maxFeeNumerator := big.NewInt(1), big.NewInt(2),
		new(big.Int).SetUint64(constants.MaxFeeNumerator)

	var tradingFeeNumerator *big.Int
	if a.Cmp(maxIndex) < 0 {
		// c + c * a
		partOne := new(big.Int).Add(
			cliffFeeNumerator,
			new(big.Int).Mul(cliffFeeNumerator, a),
		)
		// i * a * (a + 1) / 2
		partTwo := new(big.Int).Quo(
			new(big.Int).Mul(
				new(big.Int).Mul(i, a),
				new(big.Int).Add(a, one),
			),
			two,
		)
		numerator1 := new(big.Int).Add(partOne, partTwo)

		// c + i * (a + 1)
		numerator2 := new(big.Int).Add(
			cliffFeeNumerator,
			new(big.Int).Mul(
				i,
				new(big.Int).Add(a, one),
			),
		)

		firstFee, secondFee := new(big.Int).Mul(referenceAmount, numerator1),
			new(big.Int).Mul(b, numerator2)

		tradingFeeNumerator = new(big.Int).Add(firstFee, secondFee)
	} else {
		// c + (c * maxIndex)
		partOne := new(big.Int).Add(
			cliffFeeNumerator,
			new(big.Int).Mul(cliffFeeNumerator, maxIndex),
		)
		// (i * maxIndex * (maxIndex + 1)) / 2
		partTwo := new(big.Int).Quo(
			new(big.Int).Mul(
				new(big.Int).Mul(i, maxIndex),
				new(big.Int).Add(maxIndex, one),
			),
			two,
		)
		numerator1, numerator2 := new(big.Int).Add(partOne, partTwo), maxFeeNumerator

		firstFee, d := new(big.Int).Mul(referenceAmount, numerator1),
			new(big.Int).Sub(a, maxIndex)
		leftAmount := new(big.Int).Add(new(big.Int).Mul(d, referenceAmount), b)
		secondFee := new(big.Int).Mul(leftAmount, numerator2)

		tradingFeeNumerator = new(big.Int).Add(firstFee, secondFee)
	}

	denominator := new(big.Int).SetUint64(constants.FeeDenominator)
	tradingFee := new(big.Int).Div(
		new(big.Int).Sub(
			new(big.Int).Add(tradingFeeNumerator, denominator),
			one,
		),
		denominator,
	)

	// reverse to fee numerator:
	// input_amount * numerator / FEE_DENOMINATOR = trading_fee
	// => numerator = trading_fee * FEE_DENOMINATOR / input_amount
	feeNumerator, _ := MulDiv(
		tradingFee,
		denominator,
		inputAmount,
		types.RoundingUp,
	)

	if feeNumerator.Cmp(maxFeeNumerator) <= 0 {
		return feeNumerator, nil
	}

	return maxFeeNumerator, nil

}
