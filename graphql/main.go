package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("hung-han", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (cfg AppConfig) Validate() error {
	var missing []string
	if cfg.AccountURL == "" {
		missing = append(missing, "ACCOUNT_SERVICE_URL")
	}
	if cfg.CatalogURL == "" {
		missing = append(missing, "CATALOG_SERVICE_URL")
	}
	if cfg.OrderURL == "" {
		missing = append(missing, "ORDER_SERVICE_URL")
	}
	if len(missing) > 0 {
		return fmt.Errorf("set %s to run the GraphQL gateway", strings.Join(missing, ", "))
	}
	return nil
}

// 34oyRleYZT5X8Z9lrMlvl1lKblQ
