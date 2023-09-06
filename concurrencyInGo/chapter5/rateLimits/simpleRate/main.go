package main

import (
	"context"
	"log"
	"os"
	"sync"

	"golang.org/x/time/rate"
)

type APIConnection struct {
	rateLimitter *rate.Limiter
}

func Open() *APIConnection {
	// Base on token bucket
	return &APIConnection{
		// This means 1 per second token will be created , maximun 1 in total to be accomulated
		rateLimitter: rate.NewLimiter(rate.Limit(1), 1),
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
	wg.Add(20)

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
