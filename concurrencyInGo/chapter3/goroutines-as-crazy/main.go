package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	memCondumed := func() uint64 {
		// Here we invoke the garbage collector
		// This could stop the entired main thread
		runtime.GC()

		// As we clean the memory we could read the mem
		// stats by assing it to MemStats struct
		var s runtime.MemStats
		runtime.ReadMemStats(&s)

		// This is the total mem reserverd by the go
		// program, counting all ( heap, stack )
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup

	// This will unluck the wg and will be wating for data from the c channel
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e4 // 10k

	wg.Add(numGoroutines)
	before := memCondumed()

	for i := numGoroutines; i > 0; i-- {
		go noop()
	}

	wg.Wait()
	after := memCondumed()

	fmt.Printf("%.3fkb\n", float64(after-before)/numGoroutines/1000)
	fmt.Printf("%.3fMb\n", float64(after-before)/1000/1000)

}
