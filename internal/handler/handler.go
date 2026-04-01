package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/greatbody/local-service-registry/internal/checker"
	"github.com/greatbody/local-service-registry/internal/model"
	"github.com/greatbody/local-service-registry/internal/store"
)

// Handler exposes the HTTP API for the registry.
type Handler struct {
	store   *store.Store
	checker *checker.Checker
	localIP string
	mux     *http.ServeMux
}

// New creates a Handler wired to the given store and checker.
func New(s *store.Store, chk *checker.Checker) *Handler {
	h := &Handler{store: s, checker: chk, localIP: detectLocalIP()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.serveUI)
	mux.HandleFunc("GET /services", h.listServices)
	mux.HandleFunc("POST /services", h.registerService)
	mux.HandleFunc("GET /services/{id}", h.getService)
	mux.HandleFunc("DELETE /services/{id}", h.deleteService)
	h.mux = mux
	return h
}

// ServeHTTP delegates to the internal mux.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// --- handlers ---------------------------------------------------------------

func (h *Handler) serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func (h *Handler) listServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.store.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if services == nil {
		services = []*model.Service{}
	}
	for _, svc := range services {
		svc.ResolveDisplayURLs(h.localIP)
	}
	writeJSON(w, http.StatusOK, services)
}

func (h *Handler) registerService(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.URL = strings.TrimSpace(req.URL)
	if req.Name == "" || req.URL == "" {
		writeError(w, http.StatusBadRequest, "name and url are required")
		return
	}

	svc := &model.Service{
		ID:           generateID(),
		Name:         req.Name,
		URL:          req.URL,
		Description:  req.Description,
		RemoteIP:     extractIP(r),
		Status:       model.StatusUnknown,
		RegisteredAt: time.Now(),
	}

	if err := h.store.Insert(svc); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			writeError(w, http.StatusConflict, "a service with this URL is already registered")
			return
		}
		log.Printf("register: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to register service")
		return
	}

	writeJSON(w, http.StatusCreated, svc)

	// Trigger an immediate async health check for the newly registered service.
	h.checker.CheckOne(svc)
}

func (h *Handler) getService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	svc, err := h.store.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	svc.ResolveDisplayURLs(h.localIP)
	writeJSON(w, http.StatusOK, svc)
}

func (h *Handler) deleteService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ----------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
