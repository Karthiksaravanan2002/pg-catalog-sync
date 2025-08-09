package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/config"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/handler"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/local"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/repository"
	"github.com/Karthiksaravanan2002/pg-catalog-sync/internal/service"
)

func main() {
	cfg := config.Load()

	db := local.NewDB(cfg.DatabaseURL)
	defer db.Close()

	repo := repository.NewRepo(db.DB)
	svc := service.NewSyncService(repo, cfg)
	calatog := service.NewCatalogService(repo)
	h := handler.NewHandler(svc, calatog)

	h.Router()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h.Router(),
	}

	go func() {
		log.Printf("starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("server exiting")
}