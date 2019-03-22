package utils

import (
	"sync"
	"sync/atomic"
)

func increment(c int) int {
	return c + 1
}

func incrementByRef(c *int) {
	*c++
}

func incrementByRefMutex(c *int, m *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	m.Lock()
	*c++
	m.Unlock()
}

func incrementByRefAtomic(c *int32, wg *sync.WaitGroup) {
	defer wg.Done()

	atomic.AddInt32(c, 1)
}

func incrementByRefChannels(c *int, nops chan int, done chan<- bool) {
	val := <-nops
	val--
	*c++
	if val == 0 {
		done <- true
	} else {
		nops <- val
	}
}
