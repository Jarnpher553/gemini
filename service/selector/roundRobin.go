package selector

import (
	"math/rand"
	"sync"
)

func RoundRobin() Selector {
	var i = rand.Int()
	var mtx = sync.Mutex{}
	return func(nodeLens int) int {
		mtx.Lock()
		index := i % nodeLens
		i++
		mtx.Unlock()
		return index
	}
}
