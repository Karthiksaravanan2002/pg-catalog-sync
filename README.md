
# pg-catalog-sync

`pg-catalog-sync` is a Go-based microservice to sync PostgreSQL catalog metadata from an external source, persist it in a database, and expose APIs to query catalog details.

---

## Setup Instructions

### Prerequisites

- Docker & Docker Compose installed  
- Go 1.21+ (for local development)  
- PostgreSQL client (optional, for manual DB access)

### Running with Docker Compose

This project includes Dockerfiles and a `docker-compose.yml` to run all components locally:

```bash
git clone https://github.com/Karthiksaravanan2002/pg-catalog-sync.git
cd pg-catalog-sync

# Build and start all services: database, pgAdmin, mock external service, and catalog sync service
docker compose up --build
```

- PostgreSQL will be available on `localhost:5432`
- pgAdmin UI available on [http://localhost:8081](http://localhost:8081)
- Mock external metadata service on [http://localhost:8050](http://localhost:8050)
- Catalog sync service API on [http://localhost:8080](http://localhost:8080)

### Initialize Database

Run the schema SQL file to create necessary tables:

```bash
docker exec -it pgdb-1 psql -U postgres -d pgcatalog -f /path/to/schema.sql
```

> Make sure `schema.sql` is accessible inside the container, or run `psql` from the host if the DB port is mapped.

---

## API Usage Examples

### Health Check

```bash
curl http://localhost:8080/health
```

Response:

```json
{"status":"ok"}
```

---

### Sync Catalog Metadata

Send POST request to `/sync` with connection details:

```bash
curl -X POST http://localhost:8080/sync \
  -H "Content-Type: application/json" \
  -d '{
    "host": "db_host",
    "port": 5432,
    "user": "db_user",
    "password": "db_password",
    "dbname": "pgcatalog"
  }'
```

Response:

```json
{"catalog_id": "<generated_catalog_id>"}
```

---

### List All Catalogs

```bash
curl http://localhost:8080/catalogs
```

Response example:

```json
[
  {
    "id": "catalog1",
    "source": "source_info",
    "synced_at": "2025-08-09T10:00:00Z"
  }
]
```

---

### Get Catalog Details by ID

```bash
curl http://localhost:8080/catalogs/{catalog_id}
```

Response example:

```json
{
  "catalog": {
    "id": "catalog1",
    "source": "source_info",
    "synced_at": "2025-08-09T10:00:00Z"
  },
  "schemas": [
    {
      "id": 1,
      "catalog_id": "catalog1",
      "name": "public"
    }
  ]
}
```

---

## Architecture & Design Decisions

### Service Layered Design

- **Handler:** HTTP handling, request validation, and error mapping  
- **Service:** Business logic and orchestration of data flow  
- **Repository:** Database operations implemented via `sqlc` generated queries for type safety  

### Database

- PostgreSQL stores catalogs, schemas, tables, and columns metadata with foreign key constraints  
- Schema designed for referential integrity and efficient querying  

### External Metadata Service

- Mock HTTP server simulates external metadata source for development and testing  

### Context & Timeout

- Go `context.Context` with timeout ensures resiliency in DB and external service calls  

### Docker & Compose

- Multi-container Docker Compose setup enables easy local development and testing  

### SQLC for Queries

- SQL queries are defined in `.sql` files  
- `sqlc` generates type-safe Go code to avoid runtime query errors  

### Error Handling

- Clear HTTP status codes (400, 502, 504, 500) provide proper client feedback on failures  

---
