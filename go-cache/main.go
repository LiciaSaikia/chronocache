package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"chronocache/cache"

	"github.com/go-resty/resty/v2"
	_ "modernc.org/sqlite"
)

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PredictResponse struct {
	TTL int `json:"ttl"`
}

func main() {
	c, err := cache.NewChronoCache(100, 15*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	db, _ := sql.Open("sqlite", "file:cache_history.db")
	db.Exec(`CREATE TABLE IF NOT EXISTS history (
		key TEXT,
		value TEXT,
		ttl INT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		client := resty.New()
		var resp PredictResponse
		_, err := client.R().
			SetBody(map[string]string{"key": req.Key}).
			SetResult(&resp).
			SetHeader("Content-Type", "application/json").
			Post("http://127.0.0.1:8000/predict")

		if err != nil {
			http.Error(w, "Prediction failed", http.StatusInternalServerError)
			return
		}

		ttl := time.Duration(resp.TTL) * time.Second
		c.SetWithTTL(req.Key, req.Value, ttl)
		db.Exec(`INSERT INTO history (key, value, ttl) VALUES (?, ?, ?)`, req.Key, req.Value, resp.TTL)

		fmt.Fprintf(w, "Key %s set with TTL %ds", req.Key, resp.TTL)
	})

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		type EntryJSON struct {
			Key   string `json:"key"`
			Value string `json:"value"`
			TTL   int64  `json:"ttl"`
		}

		var result []EntryJSON
		for _, item := range c.Snapshot() {
			ttl := int64(item.TTLLeft.Seconds())
			if ttl < 0 {
				ttl = 0
			}
			result = append(result, EntryJSON{
				Key:   item.Key,
				Value: item.Value,
				TTL:   ttl,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		c.Delete(key)
		fmt.Fprintf(w, "Deleted key: %s", key)
	})

	log.Println("ChronoCache listening on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
