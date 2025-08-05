package messagebroker

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer interface {
	Publish(ctx context.Context, topic string, message interface{}) error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokerAddress string) Producer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokerAddress),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *kafkaProducer) Publish(ctx context.Context, topic string, message interface{}) error {
	jsonBody, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	kafkaMessage := kafka.Message{
		Topic: topic,
		Value: jsonBody,
	}

	// Menggunakan WriteMessages agar lebih tangguh
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return p.writer.WriteMessages(ctx, kafkaMessage)
}