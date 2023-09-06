package main

import (
	"context"
	"log"
	"os"
	"sync"
)

type APIConnection struct{}

func Open() *APIConnection {
	return &APIConnection{}
}

func (a *APIConnection) ReadFile(ctx context.Context) error {
	// Pretend we do work here
	return nil
}

func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// Pretend we do work here
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
