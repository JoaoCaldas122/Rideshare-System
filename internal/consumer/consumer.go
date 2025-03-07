package consumer

import (
	"context"
	"log"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func InitConsumer() {
    // Kafka consumer initialization logic
    log.Println("Kafka consumer initialized")
}

func NewConsumer(brokerAddress string, topic string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  "rideshare-group",
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) ReadMessages(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}
		c.processMessage(m)
	}
}

func (c *Consumer) processMessage(m kafka.Message) {
	// Process the message based on its key and value
	switch string(m.Key) {
	case "ride_request":
		// Handle ride request
	case "driver_location_update":
		// Handle driver location update
	default:
		log.Println("Unknown message type:", string(m.Key))
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}