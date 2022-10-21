package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	gofakeit "github.com/brianvoe/gofakeit/v6"
)

var monitoredSymbols = []string{
	"AAPL",
	"MSFT",
	"AMZN",
	"TSLA",
	"GOOGL",
	"GOOG",
	"NVDA",
	"BRK.B",
	"META",
	"UNH",
}

func main() {
	var producer sarama.AsyncProducer
	var err error

	for {
		// Update the bootstrap server address(es) below if needed. This value assumes the data
		// generator is running on the same docker network as a container named "kafka" that is
		// reachable on port 9092 for clients.
		producer, err = sarama.NewAsyncProducer([]string{"kafka:9092"}, nil)
		if err == nil {
			break
		}
		log.Printf("Got error creating producer, will retry %s", err)
		time.Sleep(1 * time.Second)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for _, s := range monitoredSymbols {
					trade := newMockTrade(s)

					bytes, err := json.Marshal(trade)
					if err != nil {
						panic(err)
					}

					log.Printf("Producing record of trade: %s", string(bytes))

					producer.Input() <- &sarama.ProducerMessage{Topic: "trades", Key: nil, Value: sarama.ByteEncoder(bytes)}
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	wg.Wait()
}

type trade struct {
	ID        int64
	Symbol    string
	Price     float64
	Size      int
	Timestamp time.Time
}

var tradeId int64

func newMockTrade(sym string) trade {
	tradeId++
	return trade{
		ID:        tradeId,
		Symbol:    sym,
		Price:     gofakeit.Price(1, 100),
		Size:      gofakeit.Number(1, 100),
		Timestamp: time.Now(),
	}
}
