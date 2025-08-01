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

type Rounding int

const (
	RoundingUp Rounding = iota
	RoundingDown
)
