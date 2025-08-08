package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Jersonmade/test-wb-project/internal/handler"
	kafkaconsumer "github.com/Jersonmade/test-wb-project/internal/kafka-consumer"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "orders",
		GroupID:   "order-consumers",
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})

	defer reader.Close()

	connStr := "postgres://postgres_user:postgres_password@localhost:5432/wb_test_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения", err)
	}

	log.Println("успешное подключение к PostgreSQL")

	defer db.Close()

	go kafkaconsumer.StartConsumer(db, reader)

	r := mux.NewRouter()
	r.HandleFunc("/orders/{orderUID}", handler.GetOrderHandler(db)).Methods("GET")

	log.Println("HTTP сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8085", r))

	
}
