package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// PriceResponse для хранения ответа от APi
type PriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// BinanceWorker реализация интерфейса Worker для работы с BinanceAPI.
type BinanceWorker struct {
	Symbols       []string
	requestsCount int
	mu            sync.Mutex
}

// Run запускает воркер, который в бесконечном цикле делает запросы к API Binance
func (w *BinanceWorker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}

	previousPrices := make(map[string]string)

	for {
		select {
		case <-ctx.Done():
			return

		default:
			for _, symbol := range w.Symbols {
				resp, err := client.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol))
				if err != nil {
					log.Println("Error fetching price:", err)
					continue
				}

				var priceResponse PriceResponse
				if err := json.NewDecoder(resp.Body).Decode(&priceResponse); err != nil {
					log.Println("Error decoding response:", err)
					resp.Body.Close()
					continue
				}
				resp.Body.Close()

				w.mu.Lock()
				w.requestsCount++
				w.mu.Unlock()

				prevPrice := previousPrices[symbol]
				if prevPrice != "" && prevPrice != priceResponse.Price {
					fmt.Printf("%s price:%s changed\n", priceResponse.Symbol, priceResponse.Price)
				} else {
					fmt.Printf("%s price:%s\n", priceResponse.Symbol, priceResponse.Price)
				}
				previousPrices[symbol] = priceResponse.Price
			}
		}
	}
}

// GetRequestsCount возвращает количество сделанных запросов
func (w *BinanceWorker) GetRequestsCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.requestsCount
}
