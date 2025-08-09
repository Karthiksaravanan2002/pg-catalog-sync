package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

type MetadataResponse struct {
	CatalogID string   `json:"catalog_id"`
	Schemas   []Schema `json:"schemas"`
}

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	resp := MetadataResponse{
		CatalogID: "abc123",
		Schemas: []Schema{
			{
				Name: "public",
				Tables: []Table{
					{
						Name: "users",
						Columns: []Column{
							{Name: "id", Type: "integer", Nullable: false},
							{Name: "email", Type: "text", Nullable: false},
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/metadata", metadataHandler)
	log.Println("Mock external service running on :8050")
	log.Fatal(http.ListenAndServe(":8050", nil))
}
