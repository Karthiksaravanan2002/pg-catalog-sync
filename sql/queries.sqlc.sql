-- name: InsertCatalog :exec
INSERT INTO catalogs (id, source, synced_at) VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET source = EXCLUDED.source, synced_at = EXCLUDED.synced_at;

-- name: InsertSchema :one
INSERT INTO schemas (catalog_id, name) VALUES ($1, $2) RETURNING id;

-- name: InsertTable :one
INSERT INTO tables (schema_id, name) VALUES ($1, $2) RETURNING id;

-- name: InsertColumn :exec
INSERT INTO columns (table_id, name, type, nullable) VALUES ($1, $2, $3, $4);

-- name: ListCatalogs :many
SELECT id, source, synced_at FROM catalogs ORDER BY synced_at DESC;

-- name: GetCatalog :one
SELECT id, source, synced_at FROM catalogs WHERE id = $1;

-- name: ListSchemasByCatalog :many
SELECT id, catalog_id, name FROM schemas WHERE catalog_id = $1;
