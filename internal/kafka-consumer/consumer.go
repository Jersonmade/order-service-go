package kafkaconsumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/Jersonmade/order-service-go/internal/model"
	"github.com/Jersonmade/order-service-go/internal/repository"
	"github.com/segmentio/kafka-go"
)

func StartConsumer(ctx context.Context, db *sql.DB, reader *kafka.Reader) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer остановлен")
			return
		default: 
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Kafka read error", err)
				continue
			}

			var order model.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Println("JSON Parse error", err)
				continue
			}

			if err := repository.SaveOrder(db, &order); err != nil {
				log.Println("Insert to Database error", err)
			} else {
				log.Printf("Order %s inserted\n", order.OrderUID)
			}
		}
	}
}