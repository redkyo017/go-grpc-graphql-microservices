package account

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "go_grpc_graphql_microservices/account/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := dialWithRetry(url)
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{
			Skip: skip,
			Take: take,
		},
	)
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, nil
}

const (
	accountDialAttempts   = 20
	accountDialTimeout    = 3 * time.Second
	accountDialRetryDelay = 2 * time.Second
)

func dialWithRetry(target string) (*grpc.ClientConn, error) {
	var lastErr error
	for attempt := 1; attempt <= accountDialAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), accountDialTimeout)
		conn, err := grpc.DialContext(
			ctx,
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()
		if err == nil {
			if attempt > 1 {
				log.Printf("account: connected to %s after %d attempts", target, attempt)
			}
			return conn, nil
		}
		lastErr = err
		log.Printf("account: dial attempt %d/%d to %s failed: %v", attempt, accountDialAttempts, target, err)
		time.Sleep(accountDialRetryDelay)
	}
	return nil, fmt.Errorf("account: unable to connect to %s after %d attempts: %w", target, accountDialAttempts, lastErr)
}
