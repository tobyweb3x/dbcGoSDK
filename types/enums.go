package types

type TokenType uint8

const (
	TokenTypeSPL TokenType = iota
	TokenTypeToken2022
)

type BaseFeeMode uint8

const (
	BaseFeeModeFeeSchedulerLinear BaseFeeMode = iota
	BaseFeeModeFeeSchedulerExponential
	BaseFeeModeFeeSchedulerRateLimiter
)

type ActivationType uint8

const (
	ActivationTypeSlot ActivationType = iota
	ActivationTypeTimestamp
)

type TradeDirection uint8

const (
	TradeDirectionBaseToQuote TradeDirection = iota
	TradeDirectionQuoteToBase
)

type CollectFeeMode uint8

const (
	CollectFeeModeQuoteToken CollectFeeMode = iota
	CollectFeeModeOutputToken
)

type Rounding uint

const (
	RoundingUp Rounding = iota
	RoundingDown
)

type MigrationOption uint8

const (
	MigrationOption_MET_DAMM    MigrationOption = 0
	MigrationOption_MET_DAMM_V2 MigrationOption = 1
)

type TokenDecimal uint8

const (
	TokenDecimal_SIX   TokenDecimal = 6
	TokenDecimal_SEVEN TokenDecimal = 7
	TokenDecimal_EIGHT TokenDecimal = 8
	TokenDecimal_NINE  TokenDecimal = 9
)

type MigrationFeeOption uint8

const (
	MigrationFeeOption_FixedBps25  MigrationFeeOption = 0
	MigrationFeeOption_FixedBps30  MigrationFeeOption = 1
	MigrationFeeOption_FixedBps100 MigrationFeeOption = 2
	MigrationFeeOption_FixedBps200 MigrationFeeOption = 3
	MigrationFeeOption_FixedBps400 MigrationFeeOption = 4
	MigrationFeeOption_FixedBps600 MigrationFeeOption = 5
)
