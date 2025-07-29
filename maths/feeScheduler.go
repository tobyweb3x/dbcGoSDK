package maths

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/types"
	"fmt"
	"math/big"
)

func GetFeeNumeratorOnLinearFeeScheduler(
	cliffFeeNumerator, reductionFactor *big.Int,
	period uint64,
) (*big.Int, error) {
	reduction := new(big.Int).Mul(
		reductionFactor, new(big.Int).SetUint64(period))

	if reduction.Cmp(cliffFeeNumerator) > 0 {
		return big.NewInt(0), nil
	}

	v := new(big.Int).Sub(cliffFeeNumerator, reduction)
	if v.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value not negative: value is %s", v.String())
	}

	return v, nil
}

func GetFeeNumeratorOnExponentialFeeScheduler(
	cliffFeeNumerator, reductionFactor *big.Int,
	period uint64,
) (*big.Int, error) {
	if period == 0 {
		return cliffFeeNumerator, nil
	}

	basisPointMax := new(big.Int).SetUint64(constants.BasisPointMax)
	if period == 1 {
		v, _ := MulDiv(
			cliffFeeNumerator,
			new(big.Int).Sub(basisPointMax, reductionFactor),
			basisPointMax,
			types.RoundingDown,
		)
		return v, nil
	}

	// calculate (1-reduction_factor/10_000)^period

	// base = ONE_Q64 - (reductionFactor << RESOLUTION) / BASIS_POINT_MAX

	reductionFactorScaled := new(big.Int).Quo(
		new(big.Int).Lsh(reductionFactor, 64),
		basisPointMax,
	)
	base := new(big.Int).Sub(constants.OneQ64, reductionFactorScaled)
	if base.Sign() < 0 {
		return nil, fmt.Errorf("safeMath requires value not negative: value is %s", base.String())
	}

	result := new(big.Int).Exp(base, new(big.Int).SetUint64(period), nil)

	return new(big.Int).Quo(
		new(big.Int).Mul(cliffFeeNumerator, result),
		constants.OneQ64,
	), nil

}
