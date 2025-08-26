package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// keyCalculator generates deterministic cache keys
type keyCalculator struct {
	includeContext    bool
	includeBuildArgs  bool
	includeTimestamps bool
}

// Calculate generates cache key from parameters
func (c *keyCalculator) Calculate(params CacheKeyParams) string {
	h := sha256.New()
	
	// Hash instruction
	h.Write([]byte(params.Instruction))
	
	// Hash context if enabled
	if c.includeContext && len(params.Context) > 0 {
		h.Write(params.Context)
	}
	
	// Hash build args in sorted order for determinism
	if c.includeBuildArgs {
		args := make([]string, 0, len(params.BuildArgs))
		for _, arg := range params.BuildArgs {
			args = append(args, arg.Key+"="+arg.Value)
		}
		sort.Strings(args)
		for _, arg := range args {
			h.Write([]byte(arg))
		}
	}
	
	// Hash base image digest
	h.Write([]byte(params.BaseImageDigest))
	
	// Hash timestamp if enabled
	if c.includeTimestamps && !params.Timestamp.IsZero() {
		h.Write([]byte(params.Timestamp.Format("2006-01-02T15:04:05Z")))
	}
	
	return hex.EncodeToString(h.Sum(nil))
}