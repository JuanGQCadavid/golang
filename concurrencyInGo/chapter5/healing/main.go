package main

import (
	"log"
	"os"
	"time"
)

type startGoroutineFn func(
	done <-chan interface{},
	pulseInterval time.Duration,
) (heartbeat <-chan interface{})

func or(channels ...<-chan interface{}) <-chan interface{} {

	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orStream := make(chan interface{})

	go func() {
		defer close(orStream)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-or(append(channels[3:], orStream)...):
			}
		}
	}()

	return orStream

}
func bridge(done <-chan interface{}, chanStream <-chan (<-chan interface{})) <-chan interface{} {

	valStream := make(chan interface{})

	go func() {
		defer close(valStream)

		for {
			stream := make(<-chan interface{})

			// Here we just pick a stream from the channel
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			// Here we are going to read values from it

			for val := range or(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}

		}
	}()

	return valStream

}

func newSteward( // This is the one that monitors the functions
	timeout time.Duration,
	startstartGoroutine startGoroutineFn,
) startGoroutineFn {
	return func(done <-chan interface{}, pulseInterval time.Duration) (heartbeat <-chan interface{}) {
		heartbeatToReturn := make(chan interface{})
		go func() {
			defer close(heartbeatToReturn)

			var wardDone chan interface{}
			var wardHearbeat <-chan interface{}

			startWard := func() { // This is the function that we are going to monitor
				wardDone = make(chan interface{})
				wardHearbeat = startstartGoroutine(or(wardDone, done), timeout/2)
			}

			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				timeout := time.After(timeout)

				for {
					select {
					case <-pulse:
						select {
						case heartbeatToReturn <- struct{}{}:
						default:
						}
					case <-wardHearbeat:
						continue monitorLoop
					case <-done:
						return
					case <-timeout:
						log.Println("Steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					}
				}
			}

		}()
		return heartbeatToReturn
	}
}
func main() {
	doWorkFn := func(
		done <-chan interface{},
		intList ...int,
	) (startGoroutineFn, <-chan interface{}) {
		intChanStream := make(chan (<-chan interface{}))
		intStream := bridge(done, intChanStream)

		doWork := func(
			done <-chan interface{},
			pulseInterval time.Duration,
		) <-chan interface{} {
			intStream := make(chan interface{})
			heartbeat := make(chan interface{})

			go func() {
				defer close(intStream)
				select {
				case intChanStream <- intStream:
				case <-done:
					return
				}
			}()

			return heartbeat
		}

		return doWork, intStream
	}
}

func simpleWard() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	doWork := func(
		done <-chan interface{},
		_ time.Duration,
	) (heartbeat <-chan interface{}) {
		log.Println("Ward: Hello, I'm irresponsible!")
		go func() {
			<-done
			log.Println("Ward: I am halting.")
		}()
		return nil
	}

	doWorkSterward := newSteward(4*time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: hlating steward and ward.")
		close(done)
	})
	for range doWorkSterward(done, 4*time.Second) {

	}
	log.Println("Done")
}
