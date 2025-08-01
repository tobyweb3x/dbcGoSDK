package maths

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/generated/dbc"
	"dbcGoSDK/types"
	"errors"
	"fmt"
	"math/big"
)

func GetFeeMode(
	collectFeeMode types.CollectFeeMode,
	tradeDirection types.TradeDirection,
	hasReferral bool,
) types.FeeMode {
	quoteToBase := tradeDirection == types.TradeDirectionQuoteToBase
	feesOnInput := quoteToBase && collectFeeMode == types.CollectFeeModeQuoteToken
	feesOnBaseToken := quoteToBase && collectFeeMode == types.CollectFeeModeOutputToken

	return types.FeeMode{
		FeeOnInput:      feesOnInput,
		FeesOnBaseToken: feesOnBaseToken,
		HasReferral:     hasReferral,
	}
}

func SwapQuote(
	virtualPool *dbc.VirtualPoolAccount,
	config *dbc.PoolConfigAccount,
	swapBaseForQuote bool,
	amountIn *big.Int,
	slippageBps uint64,
	hasReferral bool,
	currentPoint *big.Int,
) (types.QuoteResult, error) {

	if virtualPool.QuoteReserve > config.MigrationQuoteThreshold {
		return types.QuoteResult{}, errors.New("virtualPool is completed")
	}

	if amountIn.Sign() == 0 {
		return types.QuoteResult{}, errors.New("amount is zero")
	}

	tradeDirection := types.TradeDirectionQuoteToBase
	if swapBaseForQuote {
		tradeDirection = types.TradeDirectionBaseToQuote
	}

	feeMode := GetFeeMode(
		types.CollectFeeMode(config.CollectFeeMode),
		tradeDirection,
		hasReferral,
	)

	result, err := GetSwapResult(
		virtualPool,
		config,
		amountIn,
		feeMode,
		tradeDirection,
		currentPoint,
	)
	if err != nil {
		return types.QuoteResult{}, err
	}

	// calculate minimum amount out if slippage is provided
	if slippageBps > 0 {
		// slippage factor: (10000 - slippageBps) / 10000
		slippageFactor, denominator := new(big.Int).SetUint64(10_000-slippageBps),
			new(big.Int).SetUint64(10_000)

		// minimum amount out: amountOut * (10000 - slippageBps) / 10000
		minimumAmountOut := new(big.Int).Quo(
			new(big.Int).Mul(result.AmountOut, slippageFactor),
			denominator,
		)
		result.MinimumAmountOut = minimumAmountOut
		return result, nil
	}

	return result, nil
}

// GetSwapAmountFromBaseToQuote get swap amount from base to quote.
func GetSwapAmountFromBaseToQuote(
	configState []dbc.LiquidityDistributionConfig,
	currentSqrtPrice, amountIn *big.Int,
) (out types.SwapAmount, err error) {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if amountIn.Sign() == 0 {
		return types.SwapAmount{
			OutputAmount:  big.NewInt(0),
			NextSqrtPrice: currentSqrtPrice,
		}, nil
	}

	// track total output with BN
	totalOutputAmount, sqrtPrice, amountLeft :=
		big.NewInt(0), currentSqrtPrice, amountIn

	// iterate through the curve points in reverse order
	for i := len(configState) - 1; i >= 0; i-- {
		if configState[i].SqrtPrice.BigInt().Sign() == 0 ||
			configState[i].Liquidity.BigInt().Sign() == 0 {
			continue
		}

		if configState[i].SqrtPrice.BigInt().Cmp(sqrtPrice) < 0 {
			// get the current liquidity
			currentLiquidity := configState[i].Liquidity.BigInt()
			if i+1 < len(configState) {
				currentLiquidity = configState[i+1].Liquidity.BigInt()
			}

			// skip if liquidity is zero
			if currentLiquidity.Sign() == 0 {
				continue
			}

			maxAmountIn, err := GetDeltaAmountBaseUnsigned(
				configState[i].SqrtPrice.BigInt(),
				sqrtPrice,
				currentLiquidity,
				types.RoundingUp,
			)
			if err != nil {
				return types.SwapAmount{}, err
			}

			if amountLeft.Cmp(maxAmountIn) < 0 {
				nextSqrtPrice, err := GetNextSqrtPriceFromInput(
					sqrtPrice,
					currentLiquidity,
					amountLeft,
					true,
				)
				if err != nil {
					return types.SwapAmount{}, err
				}

				outputAmount, err := GetDeltaAmountQuoteUnsigned(
					nextSqrtPrice,
					sqrtPrice,
					currentLiquidity,
					types.RoundingDown,
				)
				if err != nil {
					return types.SwapAmount{}, err
				}

				totalOutputAmount.Add(totalOutputAmount, outputAmount)
				sqrtPrice, amountLeft = nextSqrtPrice, big.NewInt(0)
				break
			}

			nextSqrtPrice := configState[i].SqrtPrice.BigInt()
			outputAmount, err := GetDeltaAmountQuoteUnsigned(
				nextSqrtPrice,
				sqrtPrice,
				currentLiquidity,
				types.RoundingDown,
			)
			if err != nil {
				return types.SwapAmount{}, err
			}

			totalOutputAmount.Add(totalOutputAmount, outputAmount)
			sqrtPrice, amountLeft = nextSqrtPrice, new(big.Int).Sub(amountLeft, maxAmountIn)
		}
	}

	if amountLeft.Sign() == 0 && !(configState[0].Liquidity.BigInt().Sign() == 0) {
		nextSqrtPrice, err := GetNextSqrtPriceFromInput(
			sqrtPrice,
			configState[0].Liquidity.BigInt(),
			amountLeft,
			true,
		)
		if err != nil {
			return types.SwapAmount{}, err
		}
		outputAmount, err := GetDeltaAmountQuoteUnsigned(
			nextSqrtPrice,
			sqrtPrice,
			configState[0].Liquidity.BigInt(),
			types.RoundingDown,
		)
		if err != nil {
			return types.SwapAmount{}, err
		}

		// add to total
		totalOutputAmount.Add(totalOutputAmount, outputAmount)
		sqrtPrice = nextSqrtPrice
	}

	return types.SwapAmount{
		OutputAmount:  totalOutputAmount,
		NextSqrtPrice: sqrtPrice,
	}, nil
}

// GetInAmountFromBaseToQuote getsinput amount from base to quote (selling).
func GetInAmountFromBaseToQuote(
	configState *dbc.PoolConfigAccount,
	currentSqrtPrice, outAmount *big.Int,
) (out types.SwapAmount, err error) {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	currentSqrtPriceLocal, amountLeft, totalAmountIn :=
		currentSqrtPrice, outAmount, big.NewInt(0)

	// iterate through the curve points in reverse order
	for i := len(configState.Curve) - 1; i >= 0; i-- {
		if configState.Curve[i].SqrtPrice.BigInt().Sign() == 0 ||
			configState.Curve[i].Liquidity.BigInt().Sign() == 0 {
			continue
		}

		if configState.Curve[i].SqrtPrice.BigInt().Cmp(currentSqrtPriceLocal) < 0 {
			currentLiquidity := configState.Curve[i].Liquidity.BigInt()
			if i+1 < len(configState.Curve) {
				currentLiquidity = configState.Curve[i+1].Liquidity.BigInt()
			}

			if currentLiquidity.Sign() == 0 {
				continue
			}

			maxAmountOut, err := GetDeltaAmountQuoteUnsigned(
				configState.Curve[i].SqrtPrice.BigInt(),
				currentSqrtPriceLocal,
				currentLiquidity,
				types.RoundingDown,
			)
			if err != nil {
				return types.SwapAmount{}, err
			}

			if amountLeft.Cmp(maxAmountOut) < 0 {
				nextSqrtPrice, err := GetNextSqrtPriceFromOutput(
					currentSqrtPriceLocal,
					currentLiquidity,
					amountLeft,
					true,
				)
				if err != nil {
					return types.SwapAmount{}, err
				}

				outputAmount, err := GetDeltaAmountBaseUnsigned(
					nextSqrtPrice,
					currentSqrtPriceLocal,
					currentLiquidity,
					types.RoundingUp,
				)
				if err != nil {
					return types.SwapAmount{}, err
				}

				totalAmountIn.Add(totalAmountIn, outputAmount)
				currentSqrtPriceLocal, amountLeft = nextSqrtPrice, big.NewInt(0)
				break
			}

			nextSqrtPrice := configState.Curve[i].SqrtPrice.BigInt()
			inAmount, err := GetDeltaAmountBaseUnsigned(
				nextSqrtPrice,
				currentSqrtPriceLocal,
				currentLiquidity,
				types.RoundingUp,
			)
			if err != nil {
				return types.SwapAmount{}, err
			}

			totalAmountIn.Add(totalAmountIn, inAmount)
			currentSqrtPriceLocal, amountLeft = nextSqrtPrice, new(big.Int).Sub(amountLeft, maxAmountOut)
		}
	}

	if amountLeft.Sign() == 0 {
		nextSqrtPrice, err := GetNextSqrtPriceFromInput(
			currentSqrtPriceLocal,
			configState.Curve[0].Liquidity.BigInt(),
			amountLeft,
			true,
		)
		if err != nil {
			return types.SwapAmount{}, err
		}

		if nextSqrtPrice.Cmp(configState.SqrtStartPrice.BigInt()) < 0 {
			return types.SwapAmount{}, errors.New("not enough liquidity")
		}

		inAmount, err := GetDeltaAmountBaseUnsigned(
			nextSqrtPrice,
			currentSqrtPriceLocal,
			configState.Curve[0].Liquidity.BigInt(),
			types.RoundingUp,
		)
		if err != nil {
			return types.SwapAmount{}, err
		}

		// add to total
		totalAmountIn.Add(totalAmountIn, inAmount)
		currentSqrtPriceLocal = nextSqrtPrice
	}

	return types.SwapAmount{
		OutputAmount:  totalAmountIn,
		NextSqrtPrice: currentSqrtPriceLocal,
	}, nil
}

func GetSwapAmountFromQuoteToBase(
	configState []dbc.LiquidityDistributionConfig,
	currentSqrtPrice, amountIn *big.Int,
) (types.SwapAmount, error) {

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	if amountIn.Sign() == 0 {
		return types.SwapAmount{
			OutputAmount:  big.NewInt(0),
			NextSqrtPrice: currentSqrtPrice,
		}, nil
	}

	totalOutputAmount, sqrtPrice, amountLeft :=
		big.NewInt(0), currentSqrtPrice, amountIn

	// iterate through the curve points
	for i := range len(configState) {
		if configState[i].SqrtPrice.BigInt().Sign() == 0 ||
			configState[i].Liquidity.BigInt().Sign() == 0 {
			break
		}

		// skip if liquidity is zero
		if configState[i].Liquidity.BigInt().Sign() == 0 {
			continue
		}

		if configState[i].SqrtPrice.BigInt().Cmp(sqrtPrice) > 0 {
			// get the current liquidity
			currentLiquidity := configState[i].Liquidity.BigInt()
			maxAmountIn, err := GetDeltaAmountBaseUnsigned(
				sqrtPrice,
				configState[i].SqrtPrice.BigInt(),
				currentLiquidity,
				types.RoundingUp,
			)
			if err != nil {
				return types.SwapAmount{}, err
			}

			if amountLeft.Cmp(maxAmountIn) < 0 {
				nextSqrtPrice, err := GetNextSqrtPriceFromInput(
					sqrtPrice,
					currentLiquidity,
					amountLeft,
					true,
				)
				if err != nil {
					return types.SwapAmount{}, err
				}

				outputAmount, err := GetDeltaAmountBaseUnsigned(
					sqrtPrice,
					nextSqrtPrice,
					currentLiquidity,
					types.RoundingDown,
				)
				if err != nil {
					return types.SwapAmount{}, err
				}

				totalOutputAmount.Add(totalOutputAmount, outputAmount)
				sqrtPrice, amountLeft = nextSqrtPrice, big.NewInt(0)
				break
			}

			nextSqrtPrice := configState[i].SqrtPrice.BigInt()
			outputAmount, err := GetDeltaAmountBaseUnsigned(
				sqrtPrice,
				nextSqrtPrice,
				currentLiquidity,
				types.RoundingDown,
			)
			if err != nil {
				return types.SwapAmount{}, err
			}

			totalOutputAmount.Add(totalOutputAmount, outputAmount)
			sqrtPrice, amountLeft = nextSqrtPrice, new(big.Int).Sub(amountLeft, maxAmountIn)
		}
	}

	// check if all amount was processed
	if amountLeft.Sign() != 0 {
		return types.SwapAmount{}, errors.New("not enough liquidity to process the entire amount")
	}

	return types.SwapAmount{
		OutputAmount:  totalOutputAmount,
		NextSqrtPrice: sqrtPrice,
	}, nil
}

func GetSwapResult(
	poolState *dbc.VirtualPoolAccount,
	configState *dbc.PoolConfigAccount,
	amountIn *big.Int,
	feeMode types.FeeMode,
	tradeDirection types.TradeDirection,
	currentPoint *big.Int,
) (types.QuoteResult, error) {

	actualProtocolFee, actualTradingFee, actualReferralFee, actualAmountIn :=
		big.NewInt(0), big.NewInt(0), big.NewInt(0), amountIn

	// apply fees on input if needed
	if feeMode.FeeOnInput {
		feeResult, err := GetFeeOnAmount(
			amountIn,
			configState.PoolFees,
			feeMode.HasReferral,
			currentPoint,
			poolState.ActivationPoint,
			poolState.VolatilityTracker,
			tradeDirection,
		)
		if err != nil {
			return types.QuoteResult{}, err
		}

		actualProtocolFee = feeResult.ProtocolFee
		actualTradingFee = feeResult.TradingFee
		actualReferralFee = feeResult.ReferralFee
		actualAmountIn = feeResult.Amount
	}

	// calculate swap amount
	var (
		swapAmount types.SwapAmount
		err        error
	)
	if tradeDirection == types.TradeDirectionBaseToQuote {
		if swapAmount, err = GetSwapAmountFromBaseToQuote(
			configState.Curve[:],
			poolState.SqrtPrice.BigInt(),
			actualAmountIn,
		); err != nil {
			return types.QuoteResult{}, err
		} else {
			if swapAmount, err = GetSwapAmountFromQuoteToBase(
				configState.Curve[:],
				poolState.SqrtPrice.BigInt(),
				actualAmountIn,
			); err != nil {
				return types.QuoteResult{}, err
			}
		}
	}

	// apply fees on output if needed
	actualAmountOut := swapAmount.OutputAmount
	if !feeMode.FeeOnInput {
		feeResult, err := GetFeeOnAmount(
			swapAmount.OutputAmount,
			configState.PoolFees,
			feeMode.HasReferral,
			currentPoint,
			poolState.ActivationPoint,
			poolState.VolatilityTracker,
			tradeDirection,
		)
		if err != nil {
			return types.QuoteResult{}, err
		}
		actualProtocolFee = feeResult.ProtocolFee
		actualTradingFee = feeResult.TradingFee
		actualReferralFee = feeResult.ReferralFee
		// actualAmountIn = feeResult.Amount
	}

	return types.QuoteResult{
		AmountOut:        actualAmountOut,
		MinimumAmountOut: actualAmountOut,
		NextSqrtPrice:    swapAmount.NextSqrtPrice,
		Fee: types.QuoteFee{
			Trading:  actualTradingFee,
			Protocol: actualProtocolFee,
			Referral: actualReferralFee,
		},
		Price: types.QuotePrice{
			BeforeSwap: poolState.SqrtPrice.BigInt(),
			AfterSwap:  swapAmount.NextSqrtPrice,
		},
	}, nil
}

// CalculateQuoteExactInAmount calculates the required quote amount for exact input.
func CalculateQuoteExactInAmount(
	virtualPool *dbc.VirtualPoolAccount,
	config *dbc.PoolConfigAccount,
	currentPoint *big.Int,
) (*big.Int, error) {

	if virtualPool.QuoteReserve >= config.MigrationQuoteThreshold {
		return big.NewInt(0), nil
	}

	amountInAfterFee := config.MigrationQuoteThreshold - virtualPool.QuoteReserve

	if config.CollectFeeMode == uint8(types.CollectFeeModeQuoteToken) {
		baseFeeNumerator, err := GetBaseFeeNumerator(
			config.PoolFees.BaseFee,
			types.TradeDirectionQuoteToBase,
			currentPoint,
			new(big.Int).SetUint64(virtualPool.ActivationPoint),
			big.NewInt(0),
		)
		if err != nil {
			return nil, err
		}

		totalFeeNumerator := baseFeeNumerator
		if config.PoolFees.DynamicFee.Initialized != 0 {
			variableFee := GetVariableFee(
				config.PoolFees.DynamicFee,
				virtualPool.VolatilityTracker,
			)

			totalFeeNumerator.Add(totalFeeNumerator, variableFee)
		}

		// cap at MAX_FEE_NUMERATOR
		if hold := new(big.Int).SetUint64(constants.MaxFeeNumerator); totalFeeNumerator.Cmp(hold) >= 0 {
			totalFeeNumerator = hold
		}

		// amountIn = amountInAfterFee * FEE_DENOMINATOR / (FEE_DENOMINATOR - effectiveFeeNumerator)
		denominator := new(big.Int).Sub(new(big.Int).SetUint64(constants.FeeDenominator), totalFeeNumerator)
		return MulDiv(
			new(big.Int).SetUint64(amountInAfterFee),
			new(big.Int).SetUint64(constants.FeeDenominator),
			denominator,
			types.RoundingUp,
		)
	}

	return new(big.Int).SetUint64(amountInAfterFee), nil
}

// GetExcludedFeeAmount gets excluded fee amount from included fee amount.
func GetExcludedFeeAmount(
	tradeFeeNumerator, includedFeeAmount *big.Int,
) (struct{ ExcludedFeeAmount, TradingFee *big.Int }, error) {

	tradingFee, _ := MulDiv(
		includedFeeAmount,
		tradeFeeNumerator,
		new(big.Int).SetUint64(constants.FeeDenominator),
		types.RoundingUp,
	)

	excludedFeeAmount := new(big.Int).Sub(includedFeeAmount, tradingFee)
	if excludedFeeAmount.Sign() < 0 {
		return struct {
			ExcludedFeeAmount *big.Int
			TradingFee        *big.Int
		}{}, fmt.Errorf("safeMath requires value not negative: value is %s", excludedFeeAmount.String())
	}

	return struct {
		ExcludedFeeAmount *big.Int
		TradingFee        *big.Int
	}{
		ExcludedFeeAmount: excludedFeeAmount,
		TradingFee:        tradingFee,
	}, nil
}
func GetIncludedFeeAmount(
	tradeFeeNumerator, excludedFeeAmount *big.Int,
) (*big.Int, error) {

	includedFeeAmount, err := MulDiv(
		excludedFeeAmount,
		new(big.Int).SetUint64(constants.FeeDenominator),
		new(big.Int).Sub(
			new(big.Int).SetUint64(constants.FeeDenominator),
			tradeFeeNumerator,
		),
		types.RoundingUp,
	)
	if err != nil {
		return nil, err
	}

	// sanity check - verify the inverse calculation
	out, err := GetExcludedFeeAmount(tradeFeeNumerator, includedFeeAmount)
	if err != nil {
		return nil, err
	}

	if out.ExcludedFeeAmount.Cmp(excludedFeeAmount) < 0 {
		return nil, errors.New("inverse amount is less than excluded_fee_amount")
	}

	return includedFeeAmount, nil
}

// GetSwapResultFromOutAmount gets swap result from output amount (reverse calculation).
func GetSwapResultFromOutAmount(
	poolState *dbc.VirtualPoolAccount,
	configState *dbc.PoolConfigAccount,
	outAmount *big.Int,
	feeMode types.FeeMode,
	tradeDirection types.TradeDirection,
	currentPoint *big.Int,
) (types.QuoteResult, error) {

	actualProtocolFee, actualTradingFee, actualReferralFee :=
		big.NewInt(0), big.NewInt(0), big.NewInt(0)

	baseFeeNumerator, err := GetBaseFeeNumerator(
		configState.PoolFees.BaseFee,
		tradeDirection,
		currentPoint,
		new(big.Int).SetUint64(poolState.ActivationPoint),
		big.NewInt(0),
	)
	if err != nil {
		return types.QuoteResult{}, err
	}

	tradeFeeNumerator := baseFeeNumerator
	if configState.PoolFees.DynamicFee.Initialized != 0 {
		variableFee := GetVariableFee(
			configState.PoolFees.DynamicFee,
			poolState.VolatilityTracker,
		)
		tradeFeeNumerator.Add(tradeFeeNumerator, variableFee)
	}

	// cap at MAX_FEE_NUMERATOR
	if hold := new(big.Int).SetUint64(constants.MaxFeeNumerator); tradeFeeNumerator.Cmp(hold) >= 0 {
		tradeFeeNumerator = hold
	}

	// calculate included fee amount based on fee mode
	includedFeeOutAmount := outAmount
	if !feeMode.FeeOnInput {
		if includedFeeOutAmount, err = GetIncludedFeeAmount(tradeFeeNumerator, outAmount); err != nil {
			return types.QuoteResult{}, err
		}

		// apply fees on output if not on input
		feeResult, err := GetFeeOnAmount(
			includedFeeOutAmount,
			configState.PoolFees,
			feeMode.HasReferral,
			currentPoint,
			poolState.ActivationPoint,
			poolState.VolatilityTracker,
			tradeDirection,
		)
		if err != nil {
			return types.QuoteResult{}, err
		}

		actualProtocolFee = feeResult.ProtocolFee
		actualTradingFee = feeResult.TradingFee
		actualReferralFee = feeResult.ReferralFee
	}

	// calculate swap amount (reverse calculation)
	var swapAmount types.SwapAmount
	if tradeDirection == types.TradeDirectionBaseToQuote {
		if swapAmount, err = GetInAmountFromBaseToQuote(
			configState,
			poolState.SqrtPrice.BigInt(),
			includedFeeOutAmount,
		); err != nil {
			return types.QuoteResult{}, err
		}
	} else {
		if swapAmount, err = GetInAmountFromBaseToQuote(
			configState,
			poolState.SqrtPrice.BigInt(),
			includedFeeOutAmount,
		); err != nil {
			return types.QuoteResult{}, err
		}
	}

	// calculate included fee input amount if fees are on input
	includedFeeInAmount := swapAmount.OutputAmount
	if !feeMode.FeeOnInput {
		if includedFeeInAmount, err = GetIncludedFeeAmount(tradeFeeNumerator, swapAmount.OutputAmount); err != nil {
			return types.QuoteResult{}, err
		}

		// apply fees on input if needed
		feeResult, err := GetFeeOnAmount(
			includedFeeInAmount,
			configState.PoolFees,
			feeMode.HasReferral,
			currentPoint,
			poolState.ActivationPoint,
			poolState.VolatilityTracker,
			tradeDirection,
		)
		if err != nil {
			return types.QuoteResult{}, err
		}

		actualProtocolFee = feeResult.ProtocolFee
		actualTradingFee = feeResult.TradingFee
		actualReferralFee = feeResult.ReferralFee
	}

	return types.QuoteResult{
		AmountOut:        includedFeeInAmount,
		MinimumAmountOut: outAmount,
		NextSqrtPrice:    swapAmount.NextSqrtPrice,
		Fee: types.QuoteFee{
			Trading:  actualTradingFee,
			Protocol: actualProtocolFee,
			Referral: actualReferralFee,
		},
		Price: types.QuotePrice{
			BeforeSwap: poolState.SqrtPrice.BigInt(),
			AfterSwap:  swapAmount.NextSqrtPrice,
		},
	}, nil
}

// SwapQuoteExactOut calculates quote for a swap with exact output amount.
func SwapQuoteExactOut(
	virtualPool *dbc.VirtualPoolAccount,
	config *dbc.PoolConfigAccount,
	swapBaseForQuote bool,
	outAmount *big.Int,
	slippageBps uint64,
	hasReferral bool,
	currentPoint *big.Int,
) (types.QuoteResult, error) {

	if virtualPool.QuoteReserve >= config.MigrationQuoteThreshold {
		return types.QuoteResult{}, errors.New("virtual pool is completed")
	}

	if outAmount.Sign() == 0 {
		return types.QuoteResult{}, errors.New("amount is zero")
	}

	tradeDirection := types.TradeDirectionQuoteToBase
	if swapBaseForQuote {
		tradeDirection = types.TradeDirectionBaseToQuote
	}

	feeMode := GetFeeMode(
		types.CollectFeeMode(config.CollectFeeMode),
		tradeDirection,
		hasReferral,
	)

	result, err := GetSwapResultFromOutAmount(
		virtualPool,
		config,
		outAmount,
		feeMode,
		tradeDirection,
		currentPoint,
	)
	if err != nil {
		return types.QuoteResult{}, err
	}

	// calculate maximum amount in if slippage is provided
	if slippageBps > 0 {
		// slippage factor: (10000 + slippageBps) / 10000
		slippageFactor, denominator :=
			new(big.Int).SetUint64(10_000+slippageBps), new(big.Int).SetUint64(10_000)

		// maximum amount in: amountIn * (10000 + slippageBps) / 10000
		maximumAmountIn := new(big.Int).Div(
			new(big.Int).Mul(result.AmountOut, slippageFactor),
			denominator,
		)

		result.AmountOut = maximumAmountIn
		result.MinimumAmountOut = outAmount
		return result, nil
	}

	result.MinimumAmountOut = outAmount
	return result, nil
}
