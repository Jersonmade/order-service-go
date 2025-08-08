package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Jersonmade/test-wb-project/internal/repository"
	"github.com/gorilla/mux"
)

func GetOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderUID := vars["orderUID"]

		order, err := repository.GetOrderByUID(db, orderUID)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}
}