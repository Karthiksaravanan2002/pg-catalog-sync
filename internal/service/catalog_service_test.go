package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/domain"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/repository"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/service"
)

type mockRepo struct{}

func (m *mockRepo) InsertCatalog(ctx context.Context, id, source string, syncedAt time.Time) error {
	return nil
}

func (m *mockRepo) InsertSchema(ctx context.Context, catalogID, name string) (int64, error) {
	return 1, nil
}

func (m *mockRepo) InsertTable(ctx context.Context, schemaID int64, name string) (int64, error) {
	return 1, nil
}

func (m *mockRepo) InsertColumn(ctx context.Context, tableID int64, col domain.Column) error {
	return nil
}

func (m *mockRepo) ListCatalogs(ctx context.Context) ([]domain.Catalog, error) {
	return []domain.Catalog{
		{ID: "cat1", Source: "mocksource", SyncedAt: time.Now()},
		{ID: "cat2", Source: "mocksource", SyncedAt: time.Now()},
	}, nil
}

func (m *mockRepo) GetCatalogFull(ctx context.Context, id string) (domain.Catalog, []domain.Schema, error) {
	if id == "notfound" {
		return domain.Catalog{}, nil, repository.ErrCatalogNotFound
	}
	c := domain.Catalog{ID: id, Source: "mocksource", SyncedAt: time.Now()}
	schemas := []domain.Schema{
		{ID: 123, CatalogID: id, Name: "Schema 1"},
		{ID: 123, CatalogID: id, Name: "Schema 2"},
	}
	return c, schemas, nil
}

func TestCatalogService_ListCatalogs(t *testing.T) {
	mock := &mockRepo{}
	svc := service.NewCatalogService(mock)

	catalogs, err := svc.ListCatalogs(context.Background())
	if err != nil {
		t.Fatalf("ListCatalogs() error = %v", err)
	}
	if len(catalogs) != 2 {
		t.Fatalf("ListCatalogs() expected 2 catalogs, got %d", len(catalogs))
	}
}

func TestCatalogService_GetCatalogFull(t *testing.T) {
	mock := &mockRepo{}
	svc := service.NewCatalogService(mock)

	// Test existing catalog
	cat, schemas, err := svc.GetCatalogFull(context.Background(), "cat1")
	if err != nil {
		t.Fatalf("GetCatalogFull() error = %v", err)
	}
	if cat.ID != "cat1" {
		t.Errorf("GetCatalogFull() catalog ID = %s; want cat1", cat.ID)
	}
	if len(schemas) != 2 {
		t.Errorf("GetCatalogFull() schemas count = %d; want 2", len(schemas))
	}

	// Test not found catalog
	_, _, err = svc.GetCatalogFull(context.Background(), "notfound")
	if err != repository.ErrCatalogNotFound {
		t.Errorf("GetCatalogFull() error = %v; want ErrCatalogNotFound", err)
	}
}
