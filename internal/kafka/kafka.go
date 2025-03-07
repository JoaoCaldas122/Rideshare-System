package kafka

import (
    "context"
    "log"
    "github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func InitKafka(brokerAddress string) {
    writer = kafka.NewWriter(kafka.WriterConfig{
        Brokers:  []string{brokerAddress},
        Topic:    "ride-events",
        Balancer: &kafka.LeastBytes{},
    })
}

func PublishRideRequest(rideRequest []byte) error {
    log.Printf("Publishing ride request: %s", string(rideRequest))
    return writer.WriteMessages(context.Background(),
        kafka.Message{
            Value: rideRequest,
        },
    )
}

func PublishDriverLocationUpdate(driverLocation []byte) error {
    log.Printf("Publishing driver location update: %s", string(driverLocation))
    return writer.WriteMessages(context.Background(),
        kafka.Message{
            Value: driverLocation,
        },
    )
}

func PublishNotification(notification []byte) error {
    log.Printf("Publishing notification: %s", string(notification))
    return writer.WriteMessages(context.Background(),
        kafka.Message{
            Value: notification,
        },
    )
}

func Close() error {
    return writer.Close()
}