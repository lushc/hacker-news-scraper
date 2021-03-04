package main

import (
	"context"
	"log"
	"sync"

	"github.com/lushc/hacker-news-scraper/internal/consumer"
)

func main() {
	client := consumer.NewHNClient()
	worker := consumer.NewWorker(client)

	ctx := context.Background()
	items := make(chan int)
	wg := &sync.WaitGroup{}

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go worker.Run(ctx, items, wg)
	}

	topStories, err := client.TopStories(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range topStories {
		items <- item
	}
}
