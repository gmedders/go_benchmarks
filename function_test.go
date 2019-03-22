package utils

import (
	"sync"
	"testing"
)

func TestIncrementIntByValue(t *testing.T) {
	counter := 0
	countTo := 50000
	for i := 0; i < countTo; i++ {
		counter = increment(counter)
	}
	if counter != countTo {
		t.Errorf("Error in counting by value: %d != %d\n", counter, countTo)
	}
}

func TestIncrementIntByRef(t *testing.T) {
	counter := 0
	countTo := 50000
	for i := 0; i < countTo; i++ {
		incrementByRef(&counter)
	}
	if counter != countTo {
		t.Errorf("Error in counting by ref: %d != %d\n", counter, countTo)
	}
}

func TestIncrementIntByRefGoroutineMutex(t *testing.T) {
	counter := 0
	countTo := 50000
	m := &sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < countTo; i++ {
		wg.Add(1)
		go incrementByRefMutex(&counter, m, &wg)
	}
	wg.Wait()
	m.Lock()
	if counter != countTo {
		t.Errorf("Error in counting by ref using mutex to control goroutine: %d != %d\n", counter, countTo)
	}
	m.Unlock()
}

func TestIncrementIntByRefGoroutineAtomic(t *testing.T) {
	var counter int32
	countTo := 50000
	var wg sync.WaitGroup
	for i := 0; i < countTo; i++ {
		wg.Add(1)
		go incrementByRefAtomic(&counter, &wg)
	}
	wg.Wait()
	if counter != int32(countTo) {
		t.Errorf("Error in counting by ref using atomic to control goroutine: %d != %d\n", counter, countTo)
	}
}

func TestIncrementIntByRefGoroutineChannel(t *testing.T) {
	counter := 0
	countTo := 50000

	nops := make(chan int)
	done := make(chan bool)

	go func() {
		nops <- countTo
	}()

	for i := 0; i < countTo; i++ {
		go incrementByRefChannels(&counter, nops, done)
	}

	select {
	case <-done:
		if counter != countTo {
			t.Errorf("Error in counting by ref using channel to control goroutine: %d != %d\n", counter, countTo)
		}
	}
}

func BenchmarkIncrementIntByValue(b *testing.B) {
	counter := 0
	for i := 0; i < b.N; i++ {
		counter = increment(counter)
	}
}

func BenchmarkIncrementIntByRef(b *testing.B) {
	counter := 0
	for i := 0; i < b.N; i++ {
		incrementByRef(&counter)
	}
}

func BenchmarkIncrementIntByRefGoroutineMutex(b *testing.B) {
	counter := 0
	m := &sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go incrementByRefMutex(&counter, m, &wg)
	}
	wg.Wait()
}

func BenchmarkIncrementIntByRefGoroutineAtomic(b *testing.B) {
	var counter int32
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go incrementByRefAtomic(&counter, &wg)
	}
	wg.Wait()
}

func BenchmarkIncrementIntByRefGoroutineChannel(b *testing.B) {
	counter := 0
	countTo := b.N

	nops := make(chan int)
	done := make(chan bool)

	go func() {
		nops <- countTo
	}()

	for i := 0; i < countTo; i++ {
		go incrementByRefChannels(&counter, nops, done)
	}

	select {
	case <-done:
	}
}
