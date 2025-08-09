CREATE TABLE IF NOT EXISTS catalogs (
    id TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    synced_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS schemas (
    id BIGSERIAL PRIMARY KEY,
    catalog_id TEXT REFERENCES catalogs(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tables (
    id BIGSERIAL PRIMARY KEY,
    schema_id BIGINT REFERENCES schemas(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS columns (
    id BIGSERIAL PRIMARY KEY,
    table_id BIGINT REFERENCES tables(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    nullable BOOLEAN NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_schemas_catalog_id ON schemas (catalog_id);
CREATE INDEX IF NOT EXISTS idx_tables_schema_id ON tables (schema_id);
CREATE INDEX IF NOT EXISTS idx_columns_table_id ON columns (table_id);
