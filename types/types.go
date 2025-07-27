package types

import (
	"dbcGoSDK/generated/dbc"

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
