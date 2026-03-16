package handlers

import (
	"encoding/json"
	"net/http"

	"mpesa-finance/cache"
	"mpesa-finance/internal/database"
)

type HealthHandler struct {
	db    *database.DB
	cache *cache.RedisCache
}

func NewHealthHandler(db *database.DB, cache *cache.RedisCache) *HealthHandler {
	return &HealthHandler{
		db:    db,
		cache: cache,
	}
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health := HealthResponse{
		Status:   "healthy",
		Services: make(map[string]string),
	}
	//check database
	if err := h.db.Health(ctx); err != nil {
		health.Status = "degraded"
		health.Services["database"] = "unhealthy"
	} else {
		health.Services["database"] = "healthy"
	}
	//check redis
	if err := h.cache.Health(ctx); err != nil {
		health.Status = "degraded"
		health.Services["redis"] = "unhealthy"
	} else {
		health.Services["redis"] = "healthy"
	}
	statusCode := http.StatusOK
	if health.Status != "degraded" {
		statusCode = http.StatusServiceUnavailable
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}
