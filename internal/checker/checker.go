package checker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/greatbody/local-service-registry/internal/model"
	"github.com/greatbody/local-service-registry/internal/store"
)

// Checker periodically probes registered services.
type Checker struct {
	store    *store.Store
	interval time.Duration
	client   *http.Client
}

// New creates a Checker that runs every `interval`.
func New(s *store.Store, interval time.Duration) *Checker {
	return &Checker{
		store:    s,
		interval: interval,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Run starts the periodic health check loop. It blocks until ctx is cancelled.
func (c *Checker) Run(ctx context.Context) {
	// Run an immediate check on startup, then tick every interval.
	c.checkAll()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("health checker stopped")
			return
		case <-ticker.C:
			c.checkAll()
		}
	}
}

func (c *Checker) checkAll() {
	services, err := c.store.List()
	if err != nil {
		log.Printf("health check: failed to list services: %v", err)
		return
	}
	if len(services) == 0 {
		return
	}
	log.Printf("health check: probing %d service(s)", len(services))

	for _, svc := range services {
		c.checkOne(svc)
	}
}

// CheckOne probes a single service asynchronously and updates its status.
func (c *Checker) CheckOne(svc *model.Service) {
	go c.checkOne(svc)
}

func (c *Checker) checkOne(svc *model.Service) {
	status := c.probe(svc.URL)
	now := time.Now()
	if err := c.store.UpdateStatus(svc.ID, status, now); err != nil {
		log.Printf("health check: failed to update %s: %v", svc.Name, err)
	}
	log.Printf("  %s (%s) -> %s", svc.Name, svc.URL, status)
}

func (c *Checker) probe(url string) model.HealthStatus {
	resp, err := c.client.Get(url)
	if err != nil {
		return model.StatusUnhealthy
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return model.StatusHealthy
	}
	return model.StatusUnhealthy
}
