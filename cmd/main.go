package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jersonmade/order-service-go/internal/cache"
	"github.com/Jersonmade/order-service-go/internal/handler"
	kafkaconsumer "github.com/Jersonmade/order-service-go/internal/kafka-consumer"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "orders",
		GroupID:   "order-consumers",
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
	})

	defer reader.Close()

	connStr := "postgres://postgres_user:postgres_password@postgres:5432/wb_test_db?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("успешное подключение к PostgreSQL")

	defer db.Close()

	go kafkaconsumer.StartConsumer(ctx, db, reader)

	orderCache := cache.NewOrderCache()

	r := mux.NewRouter()
	r.HandleFunc("/orders/{orderUID}", handler.GetOrderHandler(db, orderCache)).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("HTTP сервер запущен на порту 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP сервер упал: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP graceful shutdown failed: %v", err)
	}
}
