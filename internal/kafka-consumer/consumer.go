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

func StartConsumer(db *sql.DB, reader *kafka.Reader) {
	for {
		msg, err := reader.ReadMessage(context.Background())
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
