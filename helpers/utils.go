package helpers

import (
	"fmt"
	"math/big"

	ag_binary "github.com/gagliardetto/binary"
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

func BigIntToUint128(b *big.Int) (ag_binary.Uint128, error) {
	if b.Sign() < 0 {
		return ag_binary.Uint128{}, fmt.Errorf("value must be unsigned")
	}

	if b.BitLen() > 128 {
		return ag_binary.Uint128{}, fmt.Errorf("value %s exceeds 128 bits", b.String())
	}

	var buf [16]byte
	b.FillBytes(buf[:]) // zero-pads on the left

	ag_binary.ReverseBytes(buf[:])

	var u ag_binary.Uint128
	if err := u.UnmarshalWithDecoder(ag_binary.NewBinDecoder(buf[:])); err != nil {
		return ag_binary.Uint128{}, err
	}
	return u, nil
}

// Must helper
func MustBigIntToUint128(b *big.Int) ag_binary.Uint128 {
	v, err := BigIntToUint128(b)
	if err != nil {
		panic(fmt.Errorf("cannot fit big.Int into Uint128: %s", err.Error()))
	}
	return v
}
