package wsutils

import (
	"PlantCare/websocket/wsTypes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	duration time.Duration
}

// NewRateLimiter creates and initializes a new RateLimiter
func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		duration: duration,
	}
}

// RateLimit applies rate limiting based on the event name
func (rl *RateLimiter) RateLimit(event string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Filter timestamps within the time window
	validTimes := []time.Time{}
	for _, t := range rl.requests[event] {
		if now.Sub(t) <= rl.duration {
			validTimes = append(validTimes, t)
		}
	}

	// Update the requests list for the event
	rl.requests[event] = validTimes

	// Check if the number of valid requests exceeds the limit
	if len(validTimes) >= rl.limit {
		return fmt.Errorf("rate limit exceeded for event: %s", event)
	}

	// Record the new request
	rl.requests[event] = append(rl.requests[event], now)
	return nil
}

// RateLimitWrapper wraps a handler with rate limiting logic
func (rl *RateLimiter) RateLimitWrapper(handlerFunc func(json.RawMessage, *wsTypes.Connection), event string) func(json.RawMessage, *wsTypes.Connection) {
	return func(data json.RawMessage, connection *wsTypes.Connection) {
		err := rl.RateLimit(event)
		if err != nil {
			fmt.Println(err)
			SendErrorResponse(connection, http.StatusTooManyRequests)
			return
		}
		// Call the original handler if rate limit is not exceeded
		handlerFunc(data, connection)
	}
}