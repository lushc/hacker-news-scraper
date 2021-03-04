package main

import (
	"context"
	"log"
	"sync"

	"github.com/lushc/hacker-news-scraper/internal/consumer"
	"github.com/lushc/hacker-news-scraper/internal/datastore"
)

const (
	workerCount = 10
)

func main() {
	ctx := context.Background()

	writer, err := datastore.NewDBWriter(ctx, workerCount)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	client := consumer.NewHNClient()
	worker := consumer.NewWorker(client, writer)

	items := make(chan int)
	wg := &sync.WaitGroup{}

	for i := 0; i < workerCount; i++ {
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
