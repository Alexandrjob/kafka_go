package main

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"os"
)

type Handler struct{}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf(
			"Received: topic=%s partition=%d offset=%d value=%s",
			msg.Topic,
			msg.Partition,
			msg.Offset,
			string(msg.Value),
		)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	config := sarama.NewConfig()
	config.Version = sarama.V3_5_0_0
	config.Consumer.Group.Rebalance.GroupStrategies =
		[]sarama.BalanceStrategy{
			sarama.NewBalanceStrategyRange(),
		}

	client, err := sarama.NewClient([]string{os.Getenv("KAFKA_BROKERS")}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	group, err := sarama.NewConsumerGroupFromClient(os.Getenv("GROUP_ID"), client)
	if err != nil {
		log.Fatal(err)
	}
	defer group.Close()

	handler := &Handler{}
	ctx := context.Background()

	for {
		err := group.Consume(ctx, []string{os.Getenv("TOPIC")}, handler)
		if err != nil {
			log.Printf("Consume error: %v", err)
		}
	}
}
