package helpers

import "dbcGoSDK/types"

// CheckRateLimiterApplied checks if rate limiter should be applied based on pool configuration and state.
func CheckRateLimiterApplied(
	baseFeeMode types.BaseFeeMode,
	swapBaseForQuote bool,
	currentPoint, activationPoint, maxLimiterDuration uint64,
) bool {
	return baseFeeMode == types.BaseFeeModeRateLimiter &&
		!swapBaseForQuote &&
		currentPoint >= activationPoint &&
		currentPoint <= activationPoint+maxLimiterDuration
}
