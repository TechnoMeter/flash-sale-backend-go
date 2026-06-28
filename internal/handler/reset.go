package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/TechnoMeter/FSx-flash-sale-backend-go/internal/db"
)

func ResetStock(rdb *db.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		expected := os.Getenv("RESET_KEY")
		if expected == "" {
			// If RESET_KEY is not set, deny all resets for safety.
			http.Error(w, `{"error":"reset key not configured"}`, http.StatusInternalServerError)
			return
		}
		if key != expected {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		err := rdb.Client.Set(r.Context(), "inventory:product:1", 100, 0).Err()
		if err != nil {
			http.Error(w, `{"error":"reset failed"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Stock reset to 100"})
	}
}