package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	ExternalURL string
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		log.Fatal("DATABASE_URL is required")
	}

	external := os.Getenv("EXTERNAL_METADATA_URL")
	if external == "" {
		external = "http://localhost:8050/api/metadata"
	}

	return &Config{
		Port:        port,
		DatabaseURL: db,
		ExternalURL: external,
	}
}
