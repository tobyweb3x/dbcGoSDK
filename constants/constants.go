package constants

import (
	"dbcGoSDK/generated/dammv1"
	"dbcGoSDK/generated/dammv2"
	"dbcGoSDK/generated/dbc"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

const (
	OFFSET         = 64
	RESOLUTION     = 64
	FeeDenominator = 1_000_000_000

	// MaxFeeBPS
	//  99%
	MaxFeeBPS = 9900

	// MinFeeBPS
	//  0.0001%
	MinFeeBPS = 1

	// MinFeeNumerator
	//  0.0001%
	MinFeeNumerator = 100_000

	// MaxFeeNumerator
	//  99%
	MaxFeeNumerator = 990_000_000

	// BasisPointMax
	BasisPointMax = 10000

	// MaxCurvePoint
	MaxCurvePoint = 16

	// PartnerSurplusShare
	//  80%
	PartnerSurplusShare = 80

	// SwapBufferPercentage
	//  25%
	SwapBufferPercentage = 25

	// MaxSwallowPercentage
	//  20%
	MaxSwallowPercentage = 20

	// MaxMigrationFeePercentage
	//  50%
	MaxMigrationFeePercentage = 50

	// MaxCreatorMigrationFeePercentage
	//  100%
	MaxCreatorMigrationFeePercentage = 100

	// MaxRateLimiterDurationInSeconds
	//  12 hours
	MaxRateLimiterDurationInSeconds = 43200

	// MaxRateLimiterDurationInSlots
	//  12 hours
	MaxRateLimiterDurationInSlots = 108000

	// SlotDuration
	SlotDuration = 400

	// TimestampDuration
	TimestampDuration = 1000

	// DynamicFee constants
	//

	DynamicFeeFilterPeriodDefault = 10
	DynamicFeeDecayPeriodDefault  = 120

	// DynamicFeeReductionFactorDefault
	//  50%
	DynamicFeeReductionFactorDefault = 5000

	BinStepBpsDefault = 1

	// MaxPriceChangeBpsDefault
	//  15%
	MaxPriceChangeBpsDefault = 1500
)

var (
	// Offset resolution

	// MinSqrtPrice
	//  MinSqrtPrice = big.NewInt(4295048016)
	MinSqrtPrice = big.NewInt(4295048016)

	// MaxSqrtPrice
	// MaxSqrtPrice, _ = new(big.Int).SetString("79226673521066979257578248091", 10)
	MaxSqrtPrice, _ = new(big.Int).SetString("79226673521066979257578248091", 10)

	// OneQ64 = 1 << RESOLUTION
	//  OneQ64 = new(big.Int).Lsh(big.NewInt(1), uint(RESOLUTION))
	OneQ64 = new(big.Int).Lsh(big.NewInt(1), uint(RESOLUTION))

	// DBC program ID.
	//  DBCProgramId = solana.MustPublicKeyFromBase58("dbcij3LWUppWqq96dh6gJWwBifmcGfLSB5D4DuSMaqN")
	DBCProgramId = dbc.ProgramID

	// Metaplex program ID.
	// MetaplexProgramId = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	MetaplexProgramId = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

	// DAMM v1 program ID.
	//  DammV1ProgramId = solana.MustPublicKeyFromBase58("Eo7WjKq67rjJQSZxS6z3YkapzY3eMj6Xy8X5EQVn5UaB")
	DammV1ProgramId = dammv1.ProgramID

	// DAMM v2 program ID.
	//  CpAMMProgramId = solana.MustPublicKeyFromBase58("cpamdpZCGKUy5JxQXB4dcpGPiikHawvSWAd6mEn1sGG")
	CpAMMProgramId = dammv2.ProgramID

	// Vault program ID.
	//  VaultProgramId = solana.MustPublicKeyFromBase58("24Uqj9JCLxUeoC3hGfh5W3s9FM9uCHDS2SG3LYwBpyTi")
	VaultProgramId = solana.MustPublicKeyFromBase58("24Uqj9JCLxUeoC3hGfh5W3s9FM9uCHDS2SG3LYwBpyTi")

	// Locker program ID.
	//  LockerProgramId = solana.MustPublicKeyFromBase58("LocpQgucEQHbqNABEYvBvwoxCPsSbG91A1QaQhQQqjn")
	LockerProgramId = solana.MustPublicKeyFromBase58("LocpQgucEQHbqNABEYvBvwoxCPsSbG91A1QaQhQQqjn")

	// Base address.
	//  BaseAddress = solana.MustPublicKeyFromBase58("HWzXGcGHy4tcpYfaRDCyLNzXqBTv3E6BttpCH2vJxArv")
	BaseAddress = solana.MustPublicKeyFromBase58("HWzXGcGHy4tcpYfaRDCyLNzXqBTv3E6BttpCH2vJxArv")

	// BinStepBpsU128Default
	//  BinStepBpsU128Default, _ = new(big.Int).SetString("1844674407370955", 10)
	BinStepBpsU128Default, _ = new(big.Int).SetString("1844674407370955", 10)

	// DAMM V1 Migration Fee Addresses
	//
	// DammV1MigrationFeeAddresses = []solana.PublicKey{
	// 	solana.MustPublicKeyFromBase58("8f848CEy8eY6PhJ3VcemtBDzPPSD4Vq7aJczLZ3o8MmX"),
	// 	solana.MustPublicKeyFromBase58("HBxB8Lf14Yj8pqeJ8C4qDb5ryHL7xwpuykz31BLNYr7S"),
	// 	solana.MustPublicKeyFromBase58("7v5vBdUQHTNeqk1HnduiXcgbvCyVEZ612HLmYkQoAkik"),
	// 	solana.MustPublicKeyFromBase58("EkvP7d5yKxovj884d2DwmBQbrHUWRLGK6bympzrkXGja"),
	// 	solana.MustPublicKeyFromBase58("9EZYAJrcqNWNQzP2trzZesP7XKMHA1jEomHzbRsdX8R2"),
	// 	solana.MustPublicKeyFromBase58("8cdKo87jZU2R12KY1BUjjRPwyjgdNjLGqSGQyrDshhud"),
	// }
	DammV1MigrationFeeAddresses = []solana.PublicKey{
		solana.MustPublicKeyFromBase58("8f848CEy8eY6PhJ3VcemtBDzPPSD4Vq7aJczLZ3o8MmX"),
		solana.MustPublicKeyFromBase58("HBxB8Lf14Yj8pqeJ8C4qDb5ryHL7xwpuykz31BLNYr7S"),
		solana.MustPublicKeyFromBase58("7v5vBdUQHTNeqk1HnduiXcgbvCyVEZ612HLmYkQoAkik"),
		solana.MustPublicKeyFromBase58("EkvP7d5yKxovj884d2DwmBQbrHUWRLGK6bympzrkXGja"),
		solana.MustPublicKeyFromBase58("9EZYAJrcqNWNQzP2trzZesP7XKMHA1jEomHzbRsdX8R2"),
		solana.MustPublicKeyFromBase58("8cdKo87jZU2R12KY1BUjjRPwyjgdNjLGqSGQyrDshhud"),
	}

	// DAMM V2 Migration Fee Addresses
	// DammV2MigrationFeeAddresses = []solana.PublicKey{
	// 	solana.MustPublicKeyFromBase58("7F6dnUcRuyM2TwR8myT1dYypFXpPSxqwKNSFNkxyNESd"),
	// 	solana.MustPublicKeyFromBase58("2nHK1kju6XjphBLbNxpM5XRGFj7p9U8vvNzyZiha1z6k"),
	// 	solana.MustPublicKeyFromBase58("Hv8Lmzmnju6m7kcokVKvwqz7QPmdX9XfKjJsXz8RXcjp"),
	// 	solana.MustPublicKeyFromBase58("2c4cYd4reUYVRAB9kUUkrq55VPyy2FNQ3FDL4o12JXmq"),
	// 	solana.MustPublicKeyFromBase58("AkmQWebAwFvWk55wBoCr5D62C6VVDTzi84NJuD9H7cFD"),
	// 	solana.MustPublicKeyFromBase58("DbCRBj8McvPYHJG1ukj8RE15h2dCNUdTAESG49XpQ44u"),
	// }
	DammV2MigrationFeeAddresses = []solana.PublicKey{
		solana.MustPublicKeyFromBase58("7F6dnUcRuyM2TwR8myT1dYypFXpPSxqwKNSFNkxyNESd"),
		solana.MustPublicKeyFromBase58("2nHK1kju6XjphBLbNxpM5XRGFj7p9U8vvNzyZiha1z6k"),
		solana.MustPublicKeyFromBase58("Hv8Lmzmnju6m7kcokVKvwqz7QPmdX9XfKjJsXz8RXcjp"),
		solana.MustPublicKeyFromBase58("2c4cYd4reUYVRAB9kUUkrq55VPyy2FNQ3FDL4o12JXmq"),
		solana.MustPublicKeyFromBase58("AkmQWebAwFvWk55wBoCr5D62C6VVDTzi84NJuD9H7cFD"),
		solana.MustPublicKeyFromBase58("DbCRBj8McvPYHJG1ukj8RE15h2dCNUdTAESG49XpQ44u"),
	}
)

const (
	SeedPoolAuthority           = "pool_authority"
	SeedEventAuthority          = "__event_authority"
	SeedPool                    = "pool"
	SeedTokenVault              = "token_vault"
	SeedMetadata                = "metadata"
	SeedPartnerMetadata         = "partner_metadata"
	SeedClaimFeeOperator        = "cf_operator"
	SeedDammV1MigrationMetadata = "meteora"
	SeedDammV2MigrationMetadata = "damm_v2"
	SeedLpMint                  = "lp_mint"
	SeedFee                     = "fee"
	SeedPosition                = "position"
	SeedPositionNFTAccount      = "position_nft_account"
	SeedLockEscrow              = "lock_escrow"
	SeedVirtualPoolMetadata     = "virtual_pool_metadata"
	SeedEscrow                  = "escrow"
	SeedBaseLocker              = "base_locker"
	SeedVault                   = "vault"
)
