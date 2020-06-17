package selector

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Random() Selector {

	return func(nodeLens int) int {
		i := rand.Int() % nodeLens

		return i
	}
}
