package helpers

import (
	"dbcGoSDK/types"
	"slices"

	"github.com/gagliardetto/solana-go"
)

func GetFirstkey(key1, key2 solana.PublicKey) []byte {
	if slices.Compare(key1.Bytes(), key2.Bytes()) == 1 {
		return key1.Bytes()
	}
	return key2.Bytes()
}

func GetSecondkey(key1, key2 solana.PublicKey) []byte {
	if slices.Compare(key1.Bytes(), key2.Bytes()) == 1 {
		return key2.Bytes()
	}
	return key1.Bytes()
}

// CheckRateLimiterApplied checks if rate limiter should be applied based on pool configuration and state.
func CheckRateLimiterApplied(
	baseFeeMode types.BaseFeeMode,
	swapBaseForQuote bool,
	currentPoint, activationPoint, maxLimiterDuration uint64,
) bool {
	return baseFeeMode == types.BaseFeeModeFeeSchedulerRateLimiter &&
		!swapBaseForQuote &&
		currentPoint >= activationPoint &&
		currentPoint <= activationPoint+maxLimiterDuration
}
