package service

import (
	"context"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/domain"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/repository"
)

// CatalogService provides catalog-read operations
type CatalogService struct {
	repo repository.Repository
}

func NewCatalogService(repo repository.Repository) *CatalogService {
	return &CatalogService{repo: repo}
}

// ListCatalogs returns all persisted catalogs
func (s *CatalogService) ListCatalogs(ctx context.Context) ([]domain.Catalog, error) {
	return s.repo.ListCatalogs(ctx)
}

// GetCatalogFull returns full metadata (catalog and schemas) for the given catalog id
func (s *CatalogService) GetCatalogFull(ctx context.Context, id string) (domain.Catalog, []domain.Schema, error) {
	return s.repo.GetCatalogFull(ctx, id)
}
