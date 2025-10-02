package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startTime: time.Now()}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":  "healthy",
		"message": "Server is running",
		"time":    time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) DetailedHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	response := map[string]interface{}{
		"status": "healthy",
		"uptime": uptime.String(),
		"system": map[string]interface{}{
			"goroutines":   runtime.NumGoroutine(),
			"memory_alloc": memStats.Alloc / 1024 / 1024,
			"memory_total": memStats.TotalAlloc / 1024 / 1024,
		},
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
