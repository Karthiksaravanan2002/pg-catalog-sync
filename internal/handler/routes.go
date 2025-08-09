package handler

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all HTTP routes on the provided router.
func RegisterRoutes(r *chi.Mux, h *Handler) {
	r.Get("/health", h.HealthCheck)
	r.Post("/sync", h.Sync)
	r.Get("/catalogs", h.ListCatalogs)
	r.Get("/catalogs/{id}", h.GetCatalog)
}
