package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/service"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/utils"
	"github.com/go-chi/chi/v5"
)

// Handler wraps services used by HTTP handlers.
type Handler struct {
	syncSvc    *service.SyncService
	catalogSvc *service.CatalogService
	router     *chi.Mux
}

// NewHandler constructs a Handler
func NewHandler(syncSvc *service.SyncService, catalogSvc *service.CatalogService) *Handler {
	r := chi.NewRouter()
	h := &Handler{syncSvc: syncSvc, catalogSvc: catalogSvc, router: r}
	RegisterRoutes(r, h)
	return h
}

// Router returns the underlying http.Handler (chi router)
func (h *Handler) Router() http.Handler {
	return h.router
}

// HealthCheck simple liveness endpoint
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// SyncRequest models POST /sync payload
type SyncRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"dbname"`
}

// writeJSON helper
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Sync handles POST /sync
// Validates payload, forwards it to external metadata client through service, stores metadata.
func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	var req SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON payload"})
		return
	}

	// validate required fields (password optional because we forward only)
	if !utils.IsNonEmpty(req.Host) || req.Port == 0 || !utils.IsNonEmpty(req.User) || !utils.IsNonEmpty(req.DBName) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "host, port, user and dbname are required"})
		return
	}

	// build payload map to forward to external service
	payload := map[string]interface{}{
		"host":     req.Host,
		"port":     req.Port,
		"user":     req.User,
		"password": req.Password, // forwarded in-memory only; not persisted by design
		"dbname":   req.DBName,
	}

	// timeout for the sync operation
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	catalogID, err := h.syncSvc.Sync(ctx, payload)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			writeJSON(w, http.StatusGatewayTimeout, map[string]string{"error": "external service timed out"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{"catalog_id": catalogID})
}

// ListCatalogs handles GET /catalogs
func (h *Handler) ListCatalogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	catalogs, err := h.catalogSvc.ListCatalogs(ctx)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, catalogs)
}

// GetCatalog handles GET /catalogs/{id}
func (h *Handler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
		return
	}
	catalog, schemas, err := h.catalogSvc.GetCatalogFull(ctx,id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	resp := map[string]interface{}{
		"catalog": catalog,
		"schemas": schemas,
	}
	writeJSON(w, http.StatusOK, resp)
}
