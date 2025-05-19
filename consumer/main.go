package consumer

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
	"sync"
	"time"
)

func runPoll(ctx context.Context, consumer *kafka.Consumer, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop signal received, terminate runPoll")
			return
		default:
			ev := consumer.Poll(100)
			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Message received: %s\n", string(e.Value))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				return
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "192.168.18.138:9092,192.168.18.138:9093,192.168.18.138:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "smallest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %v", err))
	}

	err = consumer.SubscribeTopics([]string{"async-topic", "sync-topic"}, nil)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	fmt.Println("Consumer initialized")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go runPoll(ctx, consumer, &wg)

	time.Sleep(time.Second * 10)
	cancel()
	wg.Wait()
}
