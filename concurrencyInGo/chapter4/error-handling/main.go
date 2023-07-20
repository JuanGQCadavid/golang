package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan interface{}, urs ...string) <-chan Result {
		results := make(chan Result)

		go func() {
			defer close(results)
			defer fmt.Println("Gorutine exited")

			for _, url := range urs {
				var result Result
				resp, err := http.Get(url)
				result = Result{
					Error:    err,
					Response: resp,
				}

				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}

	done := make(chan interface{})

	urls := []string{"https://google.com", "https://badhost", "a", "b", "c"}

	errCount := 0
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("Error: %v \n", result.Error)
			errCount++

			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}

	close(done)
	time.Sleep(1 * time.Second)
}
