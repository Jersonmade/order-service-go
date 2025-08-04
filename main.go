package main

import (
	"fmt"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=wb_user password=wb_password dbname=wb_test_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка подключения", err) 
	}

	fmt.Println("успешное подключение к PostgreSQL")

	defer db.Close()
}