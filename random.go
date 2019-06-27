package goro

import (
	"math/rand"
	"time"
)

var Random *rand.Rand

func SetSeed(seed int64) {
	source := rand.NewSource(seed)
	Random = rand.New(source)
}

func RandomSeed() int64 {
	return time.Now().UnixNano()
}
