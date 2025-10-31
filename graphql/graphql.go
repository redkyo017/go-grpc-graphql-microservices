package main

import (
	"fmt"
	"strings"

	"go_grpc_graphql_microservices/account"
	"go_grpc_graphql_microservices/catalog"
	"go_grpc_graphql_microservices/order"

	"github.com/99designs/gqlgen/graphql"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountURL, catalogURL, orderURL string) (*Server, error) {
	var missing []string
	if accountURL == "" {
		missing = append(missing, "ACCOUNT_SERVICE_URL")
	}
	if catalogURL == "" {
		missing = append(missing, "CATALOG_SERVICE_URL")
	}
	if orderURL == "" {
		missing = append(missing, "ORDER_SERVICE_URL")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing configuration: %s", strings.Join(missing, ", "))
	}

	// Connect to account service
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return nil, err
	}

	// Connect to product service
	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	// Connect to order service
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, err
	}

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
