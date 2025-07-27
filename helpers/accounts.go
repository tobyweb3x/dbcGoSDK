package helpers

import (
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
