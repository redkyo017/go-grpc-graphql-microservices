//go:generate protoc ./catalog.proto --go_out=plugins=grpc:./pb

package main

import (
	"log"
	"time"

	"go_grpc_graphql_microservices/catalog"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	r := buildRepository(cfg)
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}

func buildRepository(cfg Config) catalog.Repository {
	const (
		maxAttempts = 60
		retryDelay  = 2 * time.Second
	)

	if cfg.DatabaseURL == "" {
		log.Println("catalog: DATABASE_URL empty, using in-memory repository")
		return catalog.NewInMemoryRepository(nil)
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		repo, err := catalog.NewElasticRepository(cfg.DatabaseURL)
		if err == nil {
			log.Printf("catalog: connected to Elasticsearch at %s", cfg.DatabaseURL)
			return repo
		}

		log.Printf("catalog: elastic connection attempt %d/%d failed: %v", attempt, maxAttempts, err)
		time.Sleep(retryDelay)
	}

	log.Printf("catalog: falling back to in-memory repository after %d failed attempts", maxAttempts)
	return catalog.NewInMemoryRepository(nil)
}
