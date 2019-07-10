package goro

import (
	"math/rand"
	"time"
)

// SetSeed sets the current seed to the provided int64 seed.
func SetSeed(seed int64) {
	source := rand.NewSource(seed)
	Random = rand.New(source)
}

// RandomSeed returns a randomized int64 seed based upon time.
func RandomSeed() int64 {
	return time.Now().UnixNano()
}

// Random is our global default for random calls.
var Random = rand.New(rand.NewSource(0))
