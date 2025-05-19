package produser

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
)

func runAsyncKafka(producer *kafka.Producer) {
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	topic := "async-topic"
	message := kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}}
	for _, word := range []string{"this", "is", "asynchronous", "message", "delivery", "in", "kafka", "with", "Go", "Client"} {
		message.Value = []byte(word)
		producer.Produce(&message, nil)
	}
}

func runKafka(producer *kafka.Producer) {
	deliveryChan := make(chan kafka.Event)

	topic := "sync-topic"
	message := kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}}
	for _, word := range []string{"this", "is", "synchronous", "message", "delivery", "in", "kafka", "with", "Go", "Client"} {
		message.Value = []byte(word)
		producer.Produce(&message, deliveryChan)

		event := <-deliveryChan
		m := event.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		} else {
			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
	}

	close(deliveryChan)
}

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "192.168.18.138:9092,192.168.18.138:9093,192.168.18.138:9094",
		"acks":              "all",
		"client.id":         "myProducer",
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	defer producer.Close()
	fmt.Println("Producer initialized")

	runAsyncKafka(producer)
}
