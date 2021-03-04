package consumer

import (
	"context"
	"log"
	"sync"
)

type Worker struct {
	client Client
}

func NewWorker(client Client) *Worker {
	return &Worker{client: client}
}

func (w Worker) Run(ctx context.Context, items <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-items:
			if !ok {
				return
			}

			item, err := w.client.Item(ctx, id)
			if err != nil {
				log.Println(err)
				break
			}

			log.Println(item)
		}
	}
}
