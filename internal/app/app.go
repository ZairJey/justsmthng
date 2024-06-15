package app

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/worker"
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// Run запуск приложенияс заданными настройками
func Run(conf *config.Config) {
	maxWorkers := conf.MaxWorkers
	if maxWorkers > runtime.NumCPU() {
		maxWorkers = runtime.NumCPU()
	}

	symbolGroups := make([][]string, maxWorkers)
	for i, symbol := range conf.Symbols {
		symbolGroups[i%maxWorkers] = append(symbolGroups[i%maxWorkers], symbol)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	workers := make([]worker.Worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = worker.NewWorker(symbolGroups[i])
		wg.Add(1)
		go workers[i].Run(ctx, &wg)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				totalRequests := 0
				for _, worker := range workers {
					totalRequests += worker.GetRequestsCount()
				}
				fmt.Printf("workers requests total: %d\n", totalRequests)
			case <-ctx.Done():
				return
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "STOP" {
			cancel()
			break
		}
	}
	wg.Wait()
}
