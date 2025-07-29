package helpers

import (
	"bytes"
	"dbcGoSDK/constants"

	"github.com/gagliardetto/solana-go"
)

func DeriveDbcPoolAuthority() solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedPoolAuthority),
		},
		constants.DBCProgramId,
	)
	return pda
}

func DeriveDbcPoolMetadata(pool solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedVirtualPoolMetadata),
			pool.Bytes(),
		},
		constants.DBCProgramId,
	)
	return pda
}

func DeriveDbcPartnerMetadata(feeClaimer solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedPartnerMetadata),
			feeClaimer.Bytes(),
		},
		constants.DBCProgramId,
	)
	return pda
}

func DeriveDammV1MigrationMetadataAddress(virtualPool solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedDammV1MigrationMetadata),
			virtualPool.Bytes(),
		},
		constants.DBCProgramId,
	)
	return pda
}

// DeriveDbcTokenVaultAddress derives DBC token vault address.
func DeriveDbcTokenVaultAddress(pool, mint solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedTokenVault),
			mint.Bytes(),
			pool.Bytes(),
		},
		constants.DBCProgramId,
	)
	return pda
}

// DeriveMintMetadata derives mint metadata address.
func DeriveMintMetadata(mint solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedMetadata),
			constants.MetaplexProgramId.Bytes(),
			mint.Bytes(),
		},
		constants.MetaplexProgramId,
	)
	return pda
}

func DeriveDbcPoolAddress(quoteMint, baseMint, config solana.PublicKey) solana.PublicKey {
	isQuoteMintBiggerThanBaseMint := bytes.Compare(quoteMint.Bytes(), baseMint.Bytes()) > 0

	if isQuoteMintBiggerThanBaseMint {
		pda, _, _ := solana.FindProgramAddress(
			[][]byte{
				[]byte(constants.SeedPool),
				config.Bytes(),
				quoteMint.Bytes(),
				baseMint.Bytes(),
			},
			constants.DBCProgramId,
		)
		return pda
	}

	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte(constants.SeedPool),
			config.Bytes(),
			baseMint.Bytes(),
			quoteMint.Bytes(),
		},
		constants.DBCProgramId,
	)
	return pda
}
