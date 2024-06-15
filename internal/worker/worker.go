package worker

import (
	"awesomeProject/internal/worker/implement"
	"context"
	"sync"
)

// Worker интерфейс для воркера.
type Worker interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	GetRequestsCount() int
}

// NewWorker конструктор для создания нового воркера.
func NewWorker(symbols []string) Worker {
	return &implement.BinanceWorker{Symbols: symbols}
}
