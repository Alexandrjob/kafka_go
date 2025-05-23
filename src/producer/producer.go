package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKERS")}, config)
	if err != nil {
		log.Fatalf("Producer error: %v", err)
	}
	defer producer.Close()

	for i := 0; ; i++ {
		msg := &sarama.ProducerMessage{
			Topic: os.Getenv("TOPIC"),
			Value: sarama.StringEncoder(fmt.Sprintf("Message %d", i)),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Send error: %v", err)
			continue
		}

		log.Printf("Sent: partition=%d offset=%d", partition, offset)
	}
}
