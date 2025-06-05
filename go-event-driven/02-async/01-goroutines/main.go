package main

import (
	"log"
	"sync"
	"time"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func (h Handler) SignUp(u User) error {
	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	functions := []func(u User) error{
		h.newsletterClient.AddToNewsletter,
		h.notificationsClient.SendNotification,
	}

	// if err := h.newsletterClient.AddToNewsletter(u); err != nil {
	// 	return err
	// }

	// if err := h.notificationsClient.SendNotification(u); err != nil {
	// 	return err
	// }

	for _, f := range functions {
		wg.Add(1)
		go func(f func(u User) error) {
			log.Println("Holi")
			defer wg.Done()
			h.runFunction(u, f)

		}(f)
	}

	log.Println("Waiting")
	wg.Wait()
	log.Println("Done")

	return nil
}

func (h Handler) runFunction(u User, f func(u User) error) {
	timeToSleep := 100
	for {
		if err := f(u); err != nil {
			log.Println("Sleeping for ", timeToSleep, " seconds")
			time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
			continue
		}

		return
	}

}
