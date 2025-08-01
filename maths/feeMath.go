package maths

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/generated/dbc"
	"dbcGoSDK/helpers"
	"dbcGoSDK/types"
	"fmt"
	"math/big"
)

// GetBaseFeeNumerator get current base fee numerator.
func GetBaseFeeNumerator(
	baseFee dbc.BaseFeeConfig,
	tradeDirection types.TradeDirection,
	currentPoint *big.Int,
	activationPoint *big.Int,
	inputAmount *big.Int,
) (*big.Int, error) {

	if baseFee.BaseFeeMode == uint8(types.BaseFeeModeFeeSchedulerRateLimiter) {

		// if current point is less than activation point, return base fee
		if currentPoint.Cmp(activationPoint) < 0 {
			return new(big.Int).SetUint64(baseFee.CliffFeeNumerator), nil
		}

		maxLimiterDuration := baseFee.SecondFactor

		// if lastEffectivePoint is less than currentPoint, return base fee
		lastEffectivePoint := new(big.Int).Add(activationPoint, new(big.Int).SetUint64(maxLimiterDuration))
		if currentPoint.Cmp(lastEffectivePoint) > 0 {
			return new(big.Int).SetUint64(baseFee.CliffFeeNumerator), nil
		}

		// if no input amount provided, return base fee
		if inputAmount.Sign() == 0 {
			return new(big.Int).SetUint64(baseFee.CliffFeeNumerator), nil
		}

		// referenceAmount := baseFee.ThirdFactor
		// feeIncrementBps := baseFee.FirstFactor

		isBaseToQuote := tradeDirection == types.TradeDirectionBaseToQuote

		// check if rate limiter is applied
		isRateLimiterApplied := helpers.CheckRateLimiterApplied(
			types.BaseFeeMode(baseFee.BaseFeeMode),
			isBaseToQuote,
			currentPoint.Uint64(),
			activationPoint.Uint64(),
			maxLimiterDuration,
		)

		if isRateLimiterApplied {
			return GetFeeNumeratorOnRateLimiter(
				new(big.Int).SetUint64(baseFee.CliffFeeNumerator),
				new(big.Int).SetUint64(baseFee.ThirdFactor),
				new(big.Int).SetUint64(uint64(baseFee.FirstFactor)),
				inputAmount,
			)
		}
		return new(big.Int).SetUint64(baseFee.CliffFeeNumerator), nil
	}

	numberOfPeriod, periodFrequency, reductionFactor :=
		new(big.Int).SetUint64(uint64(baseFee.FirstFactor)), new(big.Int).SetUint64(baseFee.SecondFactor), new(big.Int).SetUint64(baseFee.ThirdFactor)

	if periodFrequency.Sign() == 0 {
		return new(big.Int).SetUint64(baseFee.CliffFeeNumerator), nil
	}

	// before activation point, use max period (min fee)
	period := numberOfPeriod
	if currentPoint.Cmp(activationPoint) >= 0 {
		elapsedPoint := new(big.Int).Sub(currentPoint, activationPoint)
		periodCount := new(big.Int).Quo(elapsedPoint, periodFrequency)

		if !periodCount.IsUint64() {
			return nil, fmt.Errorf("periodCount(%s) cannot fit into uint64", periodCount.String())
		}

		period = new(big.Int).SetUint64(min(periodCount.Uint64(), numberOfPeriod.Uint64()))
	}

	if baseFee.BaseFeeMode == uint8(types.BaseFeeModeFeeSchedulerLinear) {
		// linear fee calculation: cliffFeeNumerator - period * reductionFactor
		return GetFeeNumeratorOnLinearFeeScheduler(
			new(big.Int).SetUint64(baseFee.CliffFeeNumerator),
			reductionFactor,
			period.Uint64(),
		)
	}

	// exponential fee calculation: cliff_fee_numerator * (1 - reduction_factor/10_000)^period
	return GetFeeNumeratorOnExponentialFeeScheduler(
		new(big.Int).SetUint64(baseFee.CliffFeeNumerator),
		reductionFactor,
		period.Uint64(),
	)
}

// GetVariableFee gets variable fee from dynamic fee.
func GetVariableFee(
	dynamicFee dbc.DynamicFeeConfig,
	volatilityTracker dbc.VolatilityTracker,
) *big.Int {

	if dynamicFee.Initialized == 0 {
		return big.NewInt(0)
	}

	if volatilityTracker.VolatilityAccumulator.BigInt().Sign() == 0 {
		return big.NewInt(0)
	}

	// (volatilityAccumulator * binStep)
	volatilityTimesBinStep := new(big.Int).Mul(
		volatilityTracker.VolatilityAccumulator.BigInt(),
		new(big.Int).SetUint64(uint64(dynamicFee.BinStep)),
	)

	// (volatilityAccumulator * binStep)^2
	squared := new(big.Int).Mul(volatilityTimesBinStep, volatilityTimesBinStep)

	// (volatilityAccumulator * binStep)^2 * variableFeeControl
	vFee := new(big.Int).Mul(squared,
		new(big.Int).SetUint64(uint64(dynamicFee.VariableFeeControl)))

	scaleFactor := big.NewInt(100_000_000_000)
	numerator := new(big.Int).Add(
		vFee,
		new(big.Int).Sub(scaleFactor, big.NewInt(1)),
	)

	return new(big.Int).Quo(numerator, scaleFactor)
}

// GetFeeOnAmount get fee on amount for rate limiter.
func GetFeeOnAmount(
	amount *big.Int,
	poolFees dbc.PoolFeesConfig,
	isReferral bool,
	currentPoint *big.Int,
	activationPoint uint64,
	volatilityTracker dbc.VolatilityTracker,
	tradeDirection types.TradeDirection,
) (types.FeeOnAmountResult, error) {

	// get total trading fee
	inputAmount := big.NewInt(0)
	if poolFees.BaseFee.BaseFeeMode == uint8(types.BaseFeeModeFeeSchedulerRateLimiter) {
		inputAmount = amount
	}
	baseFeeNumerator, err := GetBaseFeeNumerator(
		poolFees.BaseFee,
		tradeDirection,
		currentPoint,
		new(big.Int).SetUint64(activationPoint),
		inputAmount,
	)
	if err != nil {
		return types.FeeOnAmountResult{}, err
	}

	// add dynamic fee if enabled
	totalFeeNumerator := baseFeeNumerator
	if poolFees.DynamicFee.Initialized != 0 {
		variableFee := GetVariableFee(
			poolFees.DynamicFee,
			volatilityTracker,
		)

		totalFeeNumerator = new(big.Int).Add(totalFeeNumerator, variableFee)
	}

	// cap at MAX_FEE_NUMERATOR
	if hold := new(big.Int).SetUint64(constants.MaxFeeNumerator); totalFeeNumerator.Cmp(hold) > 0 {
		totalFeeNumerator = hold
	}

	tradingFee, _ := MulDiv(
		amount,
		totalFeeNumerator,
		new(big.Int).SetUint64(constants.FeeDenominator),
		types.RoundingUp,
	)

	amountAfterFee := new(big.Int).Sub(amount, tradingFee)
	if amountAfterFee.Sign() < 0 {
		return types.FeeOnAmountResult{}, fmt.Errorf("safeMath requires value non-zero: value is %s", amountAfterFee.String())
	}

	protocolFee, _ := MulDiv(
		tradingFee,
		new(big.Int).SetUint64(uint64(poolFees.ProtocolFeePercent)),
		big.NewInt(100),
		types.RoundingDown,
	)

	tradingFeeAfterProtocol := new(big.Int).Sub(tradingFee, protocolFee)
	if tradingFeeAfterProtocol.Sign() < 0 {
		return types.FeeOnAmountResult{}, fmt.Errorf("safeMath requires value non-zero: value is %s", tradingFeeAfterProtocol.String())
	}

	// referral fee
	referralFee := big.NewInt(0)
	if isReferral {
		referralFee, _ = MulDiv(
			protocolFee,
			new(big.Int).SetUint64(uint64(poolFees.ProtocolFeePercent)),
			big.NewInt(100),
			types.RoundingDown,
		)
	}

	protocolFeeAfterReferral := new(big.Int).Sub(protocolFee, referralFee)
	if protocolFeeAfterReferral.Sign() < 0 {
		return types.FeeOnAmountResult{}, fmt.Errorf("safeMath requires value non-zero: value is %s", protocolFeeAfterReferral.String())
	}

	return types.FeeOnAmountResult{
		Amount:      amountAfterFee,
		TradingFee:  tradingFeeAfterProtocol,
		ProtocolFee: protocolFeeAfterReferral,
		ReferralFee: referralFee,
	}, nil

}
