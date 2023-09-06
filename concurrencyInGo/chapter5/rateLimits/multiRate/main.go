package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimit interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimit
}

func MultiLimiter(limiters ...RateLimit) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{
		limiters: limiters,
	}
}

func (l *multiLimiter) Wait(cxt context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(cxt); err != nil {
			return err
		}
	}
	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

type APIConnection struct {
	rateLimitter RateLimit
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

func Open() *APIConnection {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1) // Not burst
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)
	// Base on token bucket
	return &APIConnection{
		// This means 1 per second token will be created , maximun 1 in total to be accomulated
		rateLimitter: MultiLimiter(secondLimit, minuteLimit),
	}
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// Pretend we do work here

	// Here we are going to wait until we have the proper ammount of tokens to
	// do this operation, as it is WaitN(1) then only one token is needed

	// This could return error if the number of tokens is greatter than the burst number
	// this mean the maximoun tokens numer availabe to use
	if err := a.rateLimitter.Wait(ctx); err != nil {
		return err
	}

	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// Pretend we do work here

	if err := a.rateLimitter.Wait(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	defer log.Println("Done.")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	apiConnection := Open()
	var wg sync.WaitGroup

	time.Sleep(10 * time.Second)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())

			if err != nil {
				log.Printf("Cannot ReadFile: %v", err)
			}
			log.Println("ReadFile")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())

			if err != nil {
				log.Printf("Cannot resolveAddress: %v", err)
			}
			log.Println("ResolveAddress")
		}()
	}

	wg.Wait()
}
