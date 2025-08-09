package external_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/domain"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/external"
	"github.com/stretchr/testify/assert"
)

func TestFetchMetadata_Success(t *testing.T) {
	// Mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := domain.ExternalResponse{
			CatalogID: "abc123",
			Schemas:   []domain.ExternalSchema{},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	client := external.NewClient(mockServer.URL, 10*time.Second)
	meta, err := client.FetchMetadata(context.Background(), map[string]interface{}{})
	assert.NoError(t, err)
	assert.Equal(t, "abc123", meta.CatalogID)
}

func TestFetchMetadata_Unreachable(t *testing.T) {
	client := external.NewClient("http://127.0.0.1:9999", 10*time.Second) // Closed port
	_, err := client.FetchMetadata(context.Background(), map[string]interface{}{})
	assert.Error(t, err)
}
