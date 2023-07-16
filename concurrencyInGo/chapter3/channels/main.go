package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func main() {
	// Thet chan owner is the one that is responsbale for creating the channel
	// Sharing it with the readers, write on it an and also close it.
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)

		go func() {
			defer close(resultStream)

			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()

		return resultStream
	}

	// the Reading ownership

	resultStream := chanOwner()

	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
}

func bufferedChannelsRealCase() {
	// for this example we use the bytes.Buffer as this is a ram buffer
	// so it is faster
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	// As the channel is buffered then the channel owner
	// could write all the data fasther as it will be not blocked
	// per each read, then we could seee how the sending will be first
	// then the receving
	intStream := make(chan int, 4)

	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")

		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for interger := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v. \n", interger)
	}

}
func bufferedChannesl() {
	// Both of them are the same, a unbuffered channel that will be blocked
	// when it attemps to write to a chan but not reader is waiting
	// a := make(chan int)
	b := make(chan int, 1)

	// channels with buffers works as FIFO, so all elements will be added in order.

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		b <- 1 //As the channel has a buffer then it will add the number and exit.
		fmt.Println("Message sent to channel")

	}()

	wg.Wait()

}
func unbufferedChannels() {
	// Both of them are the same, a unbuffered channel that will be blocked
	// when it attemps to write to a chan but not reader is waiting
	a := make(chan int)
	// b := make(chan int, 0 )

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		a <- 1 // Deadlock!! as no one is readign the channel then the gorutine get blocked
		fmt.Println("Message sent to channel")

	}()

	wg.Wait()
}
func unblockRutines() {
	begin := make(chan interface{})
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin // Al gorutines will wait here until we close the channel to continue
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

func loopingOverChannel() {
	intStream := make(chan int)

	go func() {
		// This could also be faster to notify all gorutines that
		// we will not send more data, just close it, it does not matter
		// How many times the stream is being read, it will always return
		// the closed value
		defer close(intStream)

		for i := 0; i < 10; i++ {
			intStream <- i
		}
	}()

	// This range will iterate over the stream
	// until it gets closed by someone else
	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
	fmt.Println()
}

func closingChannels() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hi there"
	}()

	// The ok comes from the channel, this is a sentinel that indicates wheter the
	// answer comes from a real value or a default one as the channel was closed by the upstream
	response, ok := <-stringStream

	fmt.Printf("(%v): %s\n", ok, response)

	go func() {
		close(stringStream)
	}()

	// Ok should be false as the stream was closed
	response, ok = <-stringStream
	fmt.Printf("(%v): %s\n", ok, response)
}

func writingReading() {
	stringStream := make(chan string)

	go func() {
		// If the channnel is full then it will wait until it is empty to put the data
		// so the gorutine will be blocked here
		stringStream <- "Hi there"
	}()

	// If the channel is empty then it will block the gorotuine until some data arrive to it
	fmt.Println(<-stringStream)
}

func channels001() {
	var dataStream chan interface{} // Here we are creating a varaible of type chan, interface will be the data type that could go
	// into the channel

	dataStream = make(chan interface{}) // Here we are instantiating the channel

	// Unidirectional channels
	var onlyReadStream <-chan interface{} // This channel couuld only be used for read
	onlyReadStream = make(<-chan interface{})

	var onlyWriteStream chan<- interface{} // This channel couuld only be used for sending data
	onlyWriteStream = make(chan<- interface{})

	// Golang will adopt multi directional channels as needed into unidirectional ones

	onlyReadStream = dataStream
	onlyWriteStream = dataStream

	fmt.Println(onlyReadStream, onlyWriteStream)
}
