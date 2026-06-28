package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRaceCondition(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in CI")
	}

	// Read reset key from environment (same as in production)
	resetKey := os.Getenv("RESET_KEY")
	if resetKey == "" {
		resetKey = "reset2026" // fallback for local development
	}
	resetURL := fmt.Sprintf("http://localhost:8080/reset?key=%s", resetKey)
	resp, err := http.Get(resetURL)
	if err != nil {
		t.Fatalf("failed to reset stock: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("reset failed with status %d", resp.StatusCode)
	}

	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	successes := 0
	tooMany := 0
	mu := sync.Mutex{}

	url := "http://localhost:8080/reserve"
	totalRequests := 105 // 100 stock + 5 extra to force sold-out

	// Rate limit is 10 RPS by default; we'll send requests with 150ms delay to stay under.
	// Adjust based on your RATE_LIMIT_RPS env if changed.
	delay := 150 * time.Millisecond // > 100ms (10 RPS = 100ms between requests)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			payload := map[string]interface{}{
				"product_id": 1,
				"user_id":    fmt.Sprintf("test-user-%d", idx),
			}
			data, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")
			// No X-Test-Mode header – rate limit applies.

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()

			mu.Lock()
			if resp.StatusCode == http.StatusOK {
				successes++
			} else if resp.StatusCode == http.StatusTooManyRequests {
				tooMany++
			} else {
				t.Logf("unexpected status: %d", resp.StatusCode)
			}
			mu.Unlock()
		}(i)
		time.Sleep(delay) // stagger to avoid rate limit
	}

	wg.Wait()

	assert.Equal(t, 100, successes, "should have 100 successful reservations")
	assert.Equal(t, 5, tooMany, "should have 5 sold-out responses")
}
