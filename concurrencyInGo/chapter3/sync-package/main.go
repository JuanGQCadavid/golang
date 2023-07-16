package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup

	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)

			// fmt.Println(mem)
		}()
	}

	wg.Wait()

	fmt.Printf("%d calculators were created.", numCalcsCreated)

}

func poolGetPut() {
	// This is thread safe, so we could use get and put on gorotines with
	// confidence
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Instance")
			return struct{}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

func doOnceWithDeadlock() {
	var onceA, onceB sync.Once
	var initB func()

	initA := func() {
		onceB.Do(initB)
	}
	initB = func() {
		onceA.Do(initA)
	}

	onceA.Do(initA)
}
func doOnce() {
	// grep -ir sync.Once $(go env GOROOT)/src | wc -l
	var count int
	var once sync.Once
	var increments sync.WaitGroup

	increment := func() {
		count++
	}

	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func condBroadcast() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			// This is just a tramp to give time to the function to
			// span before the main thread ends
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()

			c.Wait()
			fn()
		}()

		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Maximizing window.")
	})

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Displaying anoing dialog box")
	})

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Mouse clicked.")
	})

	button.Clicked.Broadcast()
	clickRegistered.Wait()
}
func condSignal() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Adding to the queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}

	fmt.Println(len(queue))
}

func testMutex() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()

		count++

		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()

		count--

		fmt.Printf("Decrementing: %d\n", count)
	}

	var arithmetic sync.WaitGroup

	// Up
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}

	// Down
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("Done.")
}

func testWg() {
	hello := func(wg *sync.WaitGroup, number int) {
		// We defer it in order to ensure that when the programm finalized
		// It will tell the wg the job is done
		defer wg.Done()

		fmt.Printf("Hi %d\n", number)
	}

	greetings := 10

	var wg sync.WaitGroup
	wg.Add(greetings)

	for i := greetings; i > 0; i-- {
		go hello(&wg, i)
	}

	wg.Wait()

}
