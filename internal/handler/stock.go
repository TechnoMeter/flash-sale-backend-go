package handler

import (
	"encoding/json"
	"net/http"
	"github.com/TechnoMeter/FSx-flash-sale-backend-go/internal/db"
)

func StockHandler(rdb *db.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stock, err := rdb.Client.Get(r.Context(), "inventory:product:1").Int64()
		if err != nil {
			http.Error(w, "unavailable", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]int64{"stock": stock})
	}
}