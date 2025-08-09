package domain

import "time"

type Catalog struct {
	ID       string    `db:"id" json:"id"`
	Source   string    `db:"source" json:"source"`
	SyncedAt time.Time `db:"synced_at" json:"synced_at"`
}

type Schema struct {
	ID        int64  `db:"id" json:"id"`
	CatalogID string `db:"catalog_id" json:"catalog_id"`
	Name      string `db:"name" json:"name"`
}

type Table struct {
	ID       int64  `db:"id" json:"id"`
	SchemaID int64  `db:"schema_id" json:"schema_id"`
	Name     string `db:"name" json:"name"`
}

type Column struct {
	ID       int64  `db:"id" json:"id"`
	TableID  int64  `db:"table_id" json:"table_id"`
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`
	Nullable bool   `db:"nullable" json:"nullable"`
}

// External metadata response types
type ExternalColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type ExternalTable struct {
	Name    string           `json:"name"`
	Columns []ExternalColumn `json:"columns"`
}

type ExternalSchema struct {
	Name   string          `json:"name"`
	Tables []ExternalTable `json:"tables"`
}

type ExternalResponse struct {
	CatalogID string           `json:"catalog_id"`
	Schemas   []ExternalSchema `json:"schemas"`
}
