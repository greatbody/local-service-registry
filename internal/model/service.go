package model

import "time"

// HealthStatus represents the health state of a registered service.
type HealthStatus string

const (
	StatusUnknown   HealthStatus = "unknown"
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// Service is the core domain object stored in the registry.
type Service struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	URL           string       `json:"url"`
	Description   string       `json:"description,omitempty"`
	Status        HealthStatus `json:"status"`
	RegisteredAt  time.Time    `json:"registered_at"`
	LastCheckedAt *time.Time   `json:"last_checked_at,omitempty"`
}

// RegisterRequest is the payload for registering a new service.
type RegisterRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}
