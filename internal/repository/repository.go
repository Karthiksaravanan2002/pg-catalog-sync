package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/domain"
)

var ErrCatalogNotFound = errors.New("catalog ID not found")

type Repository interface {
    InsertCatalog(ctx context.Context, id, source string, syncedAt time.Time) error
    InsertSchema(ctx context.Context, catalogID, name string) (int64, error)
    InsertTable(ctx context.Context, schemaID int64, name string) (int64, error)
    InsertColumn(ctx context.Context, tableID int64, col domain.Column) error
    ListCatalogs(ctx context.Context) ([]domain.Catalog, error)
    GetCatalogFull(ctx context.Context, id string) (domain.Catalog, []domain.Schema, error)
}

type repo struct {
    queries *Queries
}

func NewRepo(dbConn *sql.DB) Repository {
    return &repo{
        queries: New(dbConn),
    }
}

func (r *repo) InsertCatalog(ctx context.Context, id, source string, syncedAt time.Time) error {
    return r.queries.InsertCatalog(ctx, InsertCatalogParams{
			ID: id,
			Source: source,
			SyncedAt: syncedAt,
		})
}

func (r *repo) InsertSchema(ctx context.Context, catalogID, name string) (int64, error) {
    id, err := r.queries.InsertSchema(ctx, InsertSchemaParams{
        CatalogID: catalogID,
        Name:      name,
    })
    return id, err
}

func (r *repo) InsertTable(ctx context.Context, schemaID int64, name string) (int64, error) {
    id, err := r.queries.InsertTable(ctx, InsertTableParams{
        SchemaID: schemaID,
        Name:     name,
    })
    return id, err
}

func (r *repo) InsertColumn(ctx context.Context, tableID int64, col domain.Column) error {
    return r.queries.InsertColumn(ctx, InsertColumnParams{
        TableID:  tableID,
        Name:     col.Name,
        Type:     col.Type,
        Nullable: col.Nullable,
    })
}

func (r *repo) ListCatalogs(ctx context.Context) ([]domain.Catalog, error) {
    return r.queries.ListCatalogs(ctx)
}

func (r *repo) GetCatalogFull(ctx context.Context, id string) (domain.Catalog, []domain.Schema, error) {
    c, err := r.queries.GetCatalog(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return domain.Catalog{}, nil, ErrCatalogNotFound
        }
        return domain.Catalog{}, nil, err
    }

    schemas, err := r.queries.ListSchemasByCatalog(ctx, id)
    if err != nil {
        return c, nil, err
    }
    return c, schemas, nil
}

