package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
		version:   version,
	}
}

// HealthCheck handles GET /api/health
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)

	response := map[string]interface{}{
		"status": "healthy",
		"version": h.version,
		"uptime": uptime.String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Ready handles GET /api/ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// TODO: Add readiness checks (database connection, cache, etc.)
	response := map[string]string{
		"status": "ready",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Live handles GET /api/live
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	// TODO: Add liveness checks
	response := map[string]string{
		"status": "alive",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
