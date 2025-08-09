package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/config"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/domain"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/external"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/repository"
)

type SyncService struct {
	repo repository.Repository
	ex   *external.Client
	cfg  *config.Config
}

func NewSyncService(repo repository.Repository, cfg *config.Config) *SyncService {
	return &SyncService{repo: repo, ex: external.NewClient(cfg.ExternalURL, 10*time.Second), cfg: cfg}
}

func (s *SyncService) Sync(ctx context.Context, payload map[string]interface{}) (string, error) {
	// call external service to fetch data
	rsp, err := s.ex.FetchMetadata(ctx, payload)
	if err != nil {
		return "", err
	}

	// Store Metadata
	now := time.Now().UTC()
	src := fmt.Sprintf("%s:%v/%v", payload["host"], payload["port"], payload["dbname"])
	if err := s.repo.InsertCatalog(ctx,rsp.CatalogID, src, now); err != nil {
		return "", err
	}

	// For each schema -> table -> column insert
	for _, sch := range rsp.Schemas {
		schemaID, err := s.repo.InsertSchema(ctx,rsp.CatalogID, sch.Name)
		if err != nil {
			return "", err
		}
		for _, tbl := range sch.Tables {
			tblID, err := s.repo.InsertTable(ctx,schemaID, tbl.Name)
			if err != nil {
				return "", err
			}
			for _, col := range tbl.Columns {
				c := domain.Column{Name: col.Name, Type: col.Type, Nullable: col.Nullable}
				if err := s.repo.InsertColumn(ctx,tblID, c); err != nil {
					return "", err
				}
			}
		}
	}

	return rsp.CatalogID, nil
}
