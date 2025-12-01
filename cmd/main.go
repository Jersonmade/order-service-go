package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jersonmade/order-service-go/config"
	"github.com/Jersonmade/order-service-go/internal/cache"
	"github.com/Jersonmade/order-service-go/internal/handler"
	kafkaconsumer "github.com/Jersonmade/order-service-go/internal/kafka-consumer"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   cfg.Kafka.Brokers,
		Topic:     cfg.Kafka.Topic,
		GroupID:   cfg.Kafka.GroupID,
		Partition: cfg.Kafka.Partition,
		MinBytes:  cfg.Kafka.MinBytes,
		MaxBytes:  cfg.Kafka.MaxBytes,
	})

	defer reader.Close()

	db, err := sql.Open("postgres", cfg.Database.URL)
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
		Addr:    cfg.Server.Port,
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
