package goro

import (
	"math/rand"
	"time"
)

func SetSeed(seed int64) {
	source := rand.NewSource(seed)
	Random = rand.New(source)
}

func RandomSeed() int64 {
	return time.Now().UnixNano()
}

// Ensure random is set to something!
var Random *rand.Rand = rand.New(rand.NewSource(0))
