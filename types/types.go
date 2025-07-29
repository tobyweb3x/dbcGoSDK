package types

import (
	"dbcGoSDK/generated/dbc"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type PrepareTokenAccountParams struct {
	Owner         solana.PublicKey
	Payer         solana.PublicKey
	TokenAMint    solana.PublicKey
	TokenBMint    solana.PublicKey
	TokenAProgram solana.PublicKey
	TokenBProgram solana.PublicKey
}

type FeeBreakdown struct {
	PartnerBaseFee  uint64 `json:"partnerBaseFee"`
	PartnerQuoteFee uint64 `json:"partnerQuoteFee"`
	CreatorBaseFee  uint64 `json:"creatorBaseFee"`
	CreatorQuoteFee uint64 `json:"creatorQuoteFee"`
}

type FeeTotal struct {
	TotalTradingBaseFee  uint64 `json:"totalTradingBaseFee"`
	TotalTradingQuoteFee uint64 `json:"totalTradingQuoteFee"`
}

type PoolFeeMetrics struct {
	Current FeeBreakdown `json:"current"`
	Total   FeeTotal     `json:"total"`
}

type PoolFeeByConfigOrCreator struct {
	PoolAddress          solana.PublicKey
	PartnerBaseFee       uint64
	PartnerQuoteFee      uint64
	CreatorBaseFee       uint64
	CreatorQuoteFee      uint64
	TotalTradingBaseFee  uint64
	TotalTradingQuoteFee uint64
}

type CreateVirtualPoolMetadataParam struct {
	VirtualPool solana.PublicKey
	Name        string
	Website     string
	Logo        string
	Creator     solana.PublicKey
	Payer       solana.PublicKey
}

type ClaimCreatorTradingFeeWithQuoteMintNotSolParam struct {
	Creator           solana.PublicKey
	Payer             solana.PublicKey
	FeeReceiver       solana.PublicKey
	Pool              solana.PublicKey
	PoolState         *dbc.VirtualPoolAccount
	PoolConfigState   *dbc.PoolConfigAccount
	TokenBaseProgram  solana.PublicKey
	TokenQuoteProgram solana.PublicKey
}

type ClaimCreatorTradingFeeWithQuoteMintSolParam struct {
	ClaimCreatorTradingFeeWithQuoteMintNotSolParam
	TempWSolAcc solana.PublicKey
}

type ClaimCreatorTradingFeeParam struct {
	Creator        solana.PublicKey
	Payer          solana.PublicKey
	Pool           solana.PublicKey
	MaxBaseAmount  uint64
	MaxQuoteAmount uint64
	Receiver       solana.PublicKey
	TempWSolAcc    solana.PublicKey
}

type ClaimCreatorTradingFee2Param struct {
	Creator        solana.PublicKey
	Payer          solana.PublicKey
	Pool           solana.PublicKey
	MaxBaseAmount  uint64
	MaxQuoteAmount uint64
	Receiver       solana.PublicKey
}

type CreatorWithdrawSurplusParam struct {
	Creator, VirtualPool solana.PublicKey
}

type TransferPoolCreatorParam struct {
	VirtualPool, Creator, NewCreator solana.PublicKey
}

type WithdrawMigrationFeeParam struct {
	VirtualPool solana.PublicKey
	Sender      solana.PublicKey // sender is creator or partner
	FeePayer    *solana.PublicKey
}

type InitializePoolBaseParam struct {
	Name         string
	Symbol       string
	URI          string
	Pool         solana.PublicKey
	Config       solana.PublicKey
	Payer        solana.PublicKey
	PoolCreator  solana.PublicKey
	BaseMint     solana.PublicKey
	BaseVault    solana.PublicKey
	QuoteVault   solana.PublicKey
	QuoteMint    solana.PublicKey
	MintMetadata solana.PublicKey
}

type VirtualPoolState struct {
	BaseMint solana.PublicKey
	PoolType TokenType
}

type PoolConfigState struct {
	QuoteMint      solana.PublicKey
	QuoteTokenFlag TokenType
}

type PrepareSwapParams struct {
	InputMint          solana.PublicKey
	OutputMint         solana.PublicKey
	InputTokenProgram  solana.PublicKey
	OutputTokenProgram solana.PublicKey
}

type CreatePoolParam struct {
	PreCreatePoolParam
	Payer  solana.PublicKey
	Config solana.PublicKey
}

type FirstBuyParam struct {
	Buyer                solana.PublicKey
	Receiver             solana.PublicKey
	BuyAmount            uint64
	MinimumAmountOut     uint64
	ReferralTokenAccount solana.PublicKey
}

type PreCreatePoolParam struct {
	Name        string
	Symbol      string
	URI         string
	PoolCreator solana.PublicKey
	BaseMint    solana.PublicKey
}

type CreateConfigAndPoolParam struct {
	PreCreatePoolParam PreCreatePoolParam
	CreateConfigParam
	// ConfigParameters                                       dbc.ConfigParameters
	// Config, FeeClaimer, LeftoverReceiver, QuoteMint, Payer solana.PublicKey
	TokenType TokenType
}
type CreateConfigParam struct {
	dbc.ConfigParameters
	Config, FeeClaimer, LeftoverReceiver, QuoteMint, Payer solana.PublicKey
}

type CreateConfigAndPoolWithFirstBuyParam struct {
	FirstBuyParam
	CreateConfigAndPoolParam
	BaseFeeMode
}

type CreatePoolWithFirstBuyParam struct {
	CreatePoolParam
	FirstBuyParam
}

type CreatePoolWithPartnerAndCreatorFirstBuyParam struct {
	CreatePoolParam      CreatePoolParam
	PartnerFirstBuyParam PartnerFirstBuyParam
	CreatorFirstBuyParam CreatorFirstBuyParam
}

type PartnerFirstBuyParam struct {
	Partner              solana.PublicKey
	Receiver             solana.PublicKey
	BuyAmount            uint64
	MinimumAmountOut     uint64
	ReferralTokenAccount solana.PublicKey
}
type CreatorFirstBuyParam struct {
	Creator              solana.PublicKey
	Receiver             solana.PublicKey
	BuyAmount            uint64
	MinimumAmountOut     uint64
	ReferralTokenAccount solana.PublicKey
}

type SwapParam struct {
	Owner                solana.PublicKey
	Pool                 solana.PublicKey
	AmountIn             uint64
	MinimumAmountOut     uint64
	SwapBaseForQuote     bool
	ReferralTokenAccount solana.PublicKey
	Payer                solana.PublicKey
}

type QuoteResult struct {
	AmountOut        *big.Int
	MinimumAmountOut *big.Int
	NextSqrtPrice    *big.Int
	Fee              QuoteFee
	Price            QuotePrice
}

type QuoteFee struct {
	Trading  *big.Int
	Protocol *big.Int
	Referral *big.Int // Can be nil if optional
}

type QuotePrice struct {
	BeforeSwap *big.Int
	AfterSwap  *big.Int
}

type FeeMode struct {
	FeeOnInput   bool
	FeesOnTokenA bool
	HasReferral  bool
}

type FeeOnAmountResult struct {
	Amount      *big.Int
	ProtocolFee *big.Int
	TradingFee  *big.Int
	ReferralFee *big.Int
}

type SwapAmount struct {
	OutputAmount  *big.Int
	NextSqrtPrice *big.Int
}

type SwapQuoteParam struct {
	VirtualPool      *dbc.VirtualPoolAccount
	Config           *dbc.PoolConfigAccount
	SwapBaseForQuote bool
	AmountIn         *big.Int
	SlippageBps      uint64 // optional
	HasReferral      bool
	CurrentPoint     *big.Int
}

type SwapQuoteExactInParam struct {
	VirtualPool  *dbc.VirtualPoolAccount
	Config       *dbc.PoolConfigAccount
	CurrentPoint *big.Int
}

type SwapQuoteExactOutParam struct {
	VirtualPool      *dbc.VirtualPoolAccount
	Config           *dbc.PoolConfigAccount
	SwapBaseForQuote bool
	OutAmount        *big.Int
	SlippageBps      uint64 // optional, use pointer to distinguish zero vs unset
	HasReferral      bool
	CurrentPoint     *big.Int
}

type CreatePartnerMetadataParam struct {
	Name       string
	Website    string
	Logo       string
	FeeClaimer solana.PublicKey
	Payer      solana.PublicKey
}

type ClaimPartnerTradingFeeWithQuoteMintNotSolParam struct {
	FeeClaimer        solana.PublicKey
	Payer             solana.PublicKey
	FeeReceiver       solana.PublicKey
	Config            solana.PublicKey
	Pool              solana.PublicKey
	PoolState         *dbc.VirtualPoolAccount
	PoolConfigState   *dbc.PoolConfigAccount
	TokenBaseProgram  solana.PublicKey
	TokenQuoteProgram solana.PublicKey
}

type ClaimPartnerTradingFeeWithQuoteMintSolParam struct {
	ClaimPartnerTradingFeeWithQuoteMintNotSolParam
	TempWSolAcc solana.PublicKey
}

type Accounts struct {
	PoolAuthority     solana.PublicKey
	Config            solana.PublicKey
	Pool              solana.PublicKey
	TokenAAccount     solana.PublicKey
	TokenBAccount     solana.PublicKey
	BaseVault         solana.PublicKey
	QuoteVault        solana.PublicKey
	BaseMint          solana.PublicKey
	QuoteMint         solana.PublicKey
	FeeClaimer        solana.PublicKey
	TokenBaseProgram  solana.PublicKey
	TokenQuoteProgram solana.PublicKey
}

type ClaimTradingFeeParam struct {
	FeeClaimer     solana.PublicKey
	Payer          solana.PublicKey
	Pool           solana.PublicKey
	MaxBaseAmount  *big.Int
	MaxQuoteAmount *big.Int
	Receiver       solana.PublicKey
	TempWSolAcc    solana.PublicKey
}

type PartnerWithdrawSurplusParam struct {
	FeeClaimer  solana.PublicKey
	VirtualPool solana.PublicKey
}

// type WithdrawMigrationFeeParam struct {
// 	VirtualPool solana.PublicKey
// 	Sender      solana.PublicKey  // Sender is creator or partner
// 	FeePayer    *solana.PublicKey // Optional
// }
