package worker

import (
	"context"
	"encoding/json"
	"log"
	"wallet-service/internal/service"
	pb "wallet-service/proto"

	"github.com/segmentio/kafka-go"
)

// StartConsumer memulai worker yang mendengarkan topic Kafka.
func StartConsumer(ctx context.Context, brokerAddress, topic string, walletService service.WalletService) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "wallet-service-group", // ID grup agar Kafka tahu message mana yang sudah diproses
	})

	log.Printf("Wallet Kafka consumer started on topic '%s'\n", topic)

	go func() {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// Jika context dibatalkan (aplikasi mati), hentikan loop
				if ctx.Err() != nil {
					break
				}
				log.Println("Could not read message from Kafka: ", err)
				continue
			}

			log.Printf("Received Kafka message: %s", string(m.Value))

			var event map[string]interface{}
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}
			
			// Ekstrak data dari event
			userID, _ := event["user_id"].(string)
			amount, _ := event["total_amount"].(float64)

			// Panggil logika debit di service
			_, err = walletService.Debit(context.Background(), &pb.DebitRequest{
				UserId: userID,
				Amount: amount,
			})
			
			if err != nil {
				log.Printf("Failed to process debit for user %s: %v", userID, err)
				// Di sini Anda bisa mengirim pesan ke topic 'payment_failed'
			} else {
				log.Printf("Successfully processed debit for user %s", userID)
				// Di sini Anda bisa mengirim pesan ke topic 'payment_success'
			}
		}
		r.Close()
		log.Println("Wallet Kafka consumer stopped.")
	}()
}