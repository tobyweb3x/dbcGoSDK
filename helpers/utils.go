package helpers

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func CreateProgramAccountFilter(owner solana.PublicKey, offset uint64) rpc.GetProgramAccountsOpts {
	return rpc.GetProgramAccountsOpts{
		Encoding: solana.EncodingBase58,
		Filters: []rpc.RPCFilter{
			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Bytes:  owner.Bytes(),
					Offset: offset,
				},
			},
		},
	}
}
