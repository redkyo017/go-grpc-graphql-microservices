package catalog

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "go_grpc_graphql_microservices/catalog/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := dialWithRetry(url)
	if err != nil {
		return nil, err
	}
	c := pb.NewCatalogServiceClient(conn)

	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(
		ctx,
		&pb.GetProductRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {
	r, err := c.service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Ids:   ids,
			Skip:  skip,
			Take:  take,
			Query: query,
		},
	)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, p := range r.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}

const (
	dialAttempts   = 20
	dialTimeout    = 3 * time.Second
	dialRetryDelay = 2 * time.Second
)

func dialWithRetry(target string) (*grpc.ClientConn, error) {
	var lastErr error
	for attempt := 1; attempt <= dialAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
		conn, err := grpc.DialContext(
			ctx,
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()
		if err == nil {
			if attempt > 1 {
				log.Printf("catalog: connected to %s after %d attempts", target, attempt)
			}
			return conn, nil
		}
		lastErr = err
		log.Printf("catalog: dial attempt %d/%d to %s failed: %v", attempt, dialAttempts, target, err)
		time.Sleep(dialRetryDelay)
	}
	return nil, fmt.Errorf("catalog: unable to connect to %s after %d attempts: %w", target, dialAttempts, lastErr)
}
