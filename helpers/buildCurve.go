package helpers

import (
	"dbcGoSDK/constants"
	"dbcGoSDK/generated/dbc"
	"dbcGoSDK/types"
	"fmt"
	"math"
	"math/big"
)

func BuildCurve(param types.BuildCurveParam) (dbc.ConfigParameters, error) {

	migrationBaseSupply := new(big.Float).Quo(
		new(big.Float).Mul(
			new(big.Float).SetUint64(param.TotalTokenSupply),
			new(big.Float).SetUint64(param.PercentageSupplyOnMigration),
		),
		big.NewFloat(100),
	)

	migrationQuoteAmount := GetMigrationQuoteAmountFromMigrationQuoteThreshold(
		new(big.Float).SetUint64(param.MigrationQuoteThreshold),
		param.MigrationFee.FeePercentage,
	)

	migrationPrice := new(big.Float).Quo(
		migrationQuoteAmount, migrationBaseSupply,
	)

	migrateSqrtPrice := GetSqrtPriceFromPrice(
		migrationPrice,
		param.TokenBaseDecimal,
		param.TokenQuoteDecimal,
	)

	migrationQuoteAmountInLamport := new(big.Float).Mul(
		migrationQuoteAmount,
		new(big.Float).SetFloat64(math.Pow10(int(param.TokenQuoteDecimal))),
	)
	migrationQuoteAmountInLamportBigInt, _ := migrationQuoteAmountInLamport.Int(nil)

	// ConvertToLamports(migrationQuoteAmount, param.TokenQuoteDecimal)

	migrationBaseAmount, err := GetMigrationBaseToken(
		migrationQuoteAmountInLamportBigInt,
		migrateSqrtPrice,
		param.MigrationOption,
	)
	if err != nil {
		return dbc.ConfigParameters{}, err
	}

	totalSupply := ConvertToLamports(param.TotalTokenSupply, param.TokenBaseDecimal)
	if !totalSupply.IsUint64() {
		return dbc.ConfigParameters{},
			fmt.Errorf("cannot fit totalSupply(%s) into uint64", totalSupply)
	}

	totalLeftover := ConvertToLamports(
		param.Leftover, param.TokenBaseDecimal,
	)

	lockedVesting, err := GetLockedVestingParams(
		param.LockedVestingParam.TotalLockedVestingAmount,
		param.LockedVestingParam.NumberOfVestingPeriod,
		param.LockedVestingParam.CliffUnlockAmount,
		param.LockedVestingParam.TotalVestingDuration,
		param.LockedVestingParam.CliffDurationFromMigrationTime,
		param.TokenBaseDecimal,
	)
	if err != nil {
		return dbc.ConfigParameters{}, err
	}

	totalVestingAmount := GetTotalVestingAmount(
		dbc.LockedVestingParams{
			AmountPerPeriod:                lockedVesting.AmountPerPeriod,
			CliffDurationFromMigrationTime: lockedVesting.CliffDurationFromMigrationTime,
			Frequency:                      lockedVesting.Frequency,
			NumberOfPeriod:                 lockedVesting.NumberOfPeriod,
			CliffUnlockAmount:              lockedVesting.CliffUnlockAmount,
		},
	)

	swapAmount := new(big.Int).Sub(
		new(big.Int).Sub(
			new(big.Int).Sub(totalSupply, migrationBaseAmount),
			totalVestingAmount,
		),
		totalLeftover,
	)

	migrationQuoteThresholdInLamport := ConvertToLamports(
		param.MigrationQuoteThreshold, param.TokenQuoteDecimal,
	)

	if !migrationQuoteThresholdInLamport.IsUint64() {
		return dbc.ConfigParameters{},
			fmt.Errorf("cannot fit migrationQuoteThresholdInLamport(%s) as uint64", migrationQuoteThresholdInLamport)
	}

	firstCurve, err := GetFirstCurve(
		migrateSqrtPrice,
		migrationBaseAmount,
		swapAmount,
		migrationQuoteThresholdInLamport,
		param.MigrationFee.FeePercentage,
	)
	if err != nil {
		return dbc.ConfigParameters{}, err
	}

	totalDynamicSupply, err := GetTotalSupplyFromCurve(
		migrationQuoteThresholdInLamport,
		firstCurve.SqrtStartPrice,
		firstCurve.Curve,
		lockedVesting,
		param.MigrationOption,
		totalLeftover,
		param.MigrationFee.FeePercentage,
	)
	if err != nil {
		return dbc.ConfigParameters{}, err
	}

	remainingAmount := new(big.Int).Sub(totalSupply, totalDynamicSupply)

	lastLiquidity, err := GetInitialLiquidityFromDeltaBase(
		remainingAmount,
		constants.MaxSqrtPrice,
		migrateSqrtPrice,
	)
	if err != nil {
		return dbc.ConfigParameters{}, err
	}

	if lastLiquidity.Sign() != 0 {
		firstCurve.Curve = append(firstCurve.Curve, dbc.LiquidityDistributionParameters{
			SqrtPrice: MustBigIntToUint128(constants.MaxSqrtPrice),
			Liquidity: MustBigIntToUint128(lastLiquidity),
		})
	}

	baseFee, err := GetBaseFeeParams(
		param.BaseFeeParams,
		param.TokenQuoteDecimal,
		param.ActivationType,
	)
	if err != nil {
		return dbc.ConfigParameters{}, err
	}

	var dynamicField *dbc.DynamicFeeParameters
	if param.DynamicFeeEnabled {
		baseFeeBp := param.BaseFeeParams.FeeSchedulerParam.EndingFeeBps
		if param.BaseFeeParams.BaseFeeMode == types.BaseFeeModeFeeSchedulerRateLimiter {
			baseFeeBp = param.BaseFeeParams.RateLimiterParam.BaseFeeBps
		}
		d, err := GetDynamicFeeParams(baseFeeBp, 0)
		if err != nil {
			return dbc.ConfigParameters{}, err
		}
		dynamicField = &d
	}

	return dbc.ConfigParameters{
		PoolFees: dbc.PoolFeeParameters{
			BaseFee: dbc.BaseFeeParameters{
				CliffFeeNumerator: baseFee.CliffFeeNumerator,
				FirstFactor:       baseFee.FirstFactor,
				SecondFactor:      baseFee.SecondFactor,
				ThirdFactor:       baseFee.ThirdFactor,
				BaseFeeMode:       uint8(baseFee.BaseFeeMode),
			},
			DynamicFee: dynamicField,
		},
		ActivationType:            uint8(param.ActivationType),
		CollectFeeMode:            uint8(param.CollectFeeMode),
		MigrationOption:           uint8(param.MigrationFeeOption),
		TokenType:                 uint8(param.TokenType),
		TokenDecimal:              uint8(param.TokenBaseDecimal),
		MigrationQuoteThreshold:   migrationQuoteThresholdInLamport.Uint64(),
		PartnerLpPercentage:       param.PartnerLpPercentage,
		CreatorLpPercentage:       param.CreatorLpPercentage,
		PartnerLockedLpPercentage: param.PartnerLockedLpPercentage,
		CreatorLockedLpPercentage: param.CreatorLockedLpPercentage,
		SqrtStartPrice:            MustBigIntToUint128(firstCurve.SqrtStartPrice),
		LockedVesting:             lockedVesting,
		MigrationFeeOption:        uint8(param.MigrationFeeOption),
		TokenSupply: &dbc.TokenSupplyParams{
			PreMigrationTokenSupply:  totalSupply.Uint64(),
			PostMigrationTokenSupply: totalSupply.Uint64(),
		},
		CreatorTradingFeePercentage: param.CreatorTradingFeePercentage,
		TokenUpdateAuthority:        param.TokenUpdateAuthority,
		MigrationFee:                param.MigrationFee,
		Curve:                       firstCurve.Curve,
	}, nil

}

func GetDynamicFeeParams(
	baseFeeBp, maxPriceChangeBps uint64,
) (dbc.DynamicFeeParameters, error) {
	if maxPriceChangeBps > constants.MaxPriceChangeBpsDefault {
		return dbc.DynamicFeeParameters{},
			fmt.Errorf("maxPriceChangeBps (%d bps) must be less than or equal to %d",
				maxPriceChangeBps, constants.MaxPriceChangeBpsDefault)
	}

	if maxPriceChangeBps == 0 {
		maxPriceChangeBps = constants.MaxPriceChangeBpsDefault
	}
	priceRatio := maxPriceChangeBps/constants.BasisPointMax + 1

	hold := new(big.Float).SetUint64(priceRatio)
	hold.Sqrt(hold).Mul(hold, big.NewFloat(math.Pow(2, 64)))

	sqrtPriceRatioQ64 := new(big.Int)
	hold.Int(sqrtPriceRatioQ64)

	deltaBinId := new(big.Int).Mul(
		new(big.Int).Quo(
			new(big.Int).Sub(sqrtPriceRatioQ64, constants.OneQ64),
			constants.BinStepBpsU128Default,
		),
		big.NewInt(2),
	)

	maxVolatilityAccumulator := deltaBinId.Mul(deltaBinId, big.NewInt(constants.BasisPointMax))

	squareVfaBin := new(big.Int).Mul(maxVolatilityAccumulator, constants.BinStepBpsU128Default)
	squareVfaBin.Exp(squareVfaBin, big.NewInt(2), nil)

	baseFeeNumerator := BpsToFeeNumerator(baseFeeBp)
	maxDynamicFeeNumerator := baseFeeNumerator.Mul(baseFeeNumerator, big.NewInt(20))
	maxDynamicFeeNumerator.Quo(maxDynamicFeeNumerator, big.NewInt(100)) // default max dynamic fee = 20% of min base fee

	vFee := maxDynamicFeeNumerator.Mul(maxDynamicFeeNumerator, big.NewInt(100_000_000_000))
	vFee.Sub(vFee, big.NewInt(99_999_999_999))

	variableFeeControl := vFee.Quo(vFee, squareVfaBin)

	if maxVolatilityAccumulator.Cmp(big.NewInt(math.MaxUint32)) > 0 ||
		variableFeeControl.Cmp(big.NewInt(math.MaxUint32)) > 0 {
		return dbc.DynamicFeeParameters{}, fmt.Errorf("either variableFeeControl(%s) or maxVolatilityAccumulator(%s) cannot fit into uint32",
			variableFeeControl, maxVolatilityAccumulator)
	}

	return dbc.DynamicFeeParameters{
		BinStep:                  constants.BinStepBpsDefault,
		BinStepU128:              MustBigIntToUint128(constants.BinStepBpsU128Default),
		FilterPeriod:             constants.DynamicFeeFilterPeriodDefault,
		DecayPeriod:              constants.DynamicFeeDecayPeriodDefault,
		ReductionFactor:          constants.DynamicFeeReductionFactorDefault,
		MaxVolatilityAccumulator: uint32(maxVolatilityAccumulator.Uint64()),
		VariableFeeControl:       uint32(variableFeeControl.Uint64()),
	}, nil
}
