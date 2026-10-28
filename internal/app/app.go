package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)


func NewPool(ctx context.Context) (*pgxpool.Pool, error){
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil{
		return nil,err
	}

	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	pool,err := pgxpool.NewWithConfig(ctx, config)
	if err != nil{
		pool.Close()
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil{
		pool.Close()
		return nil, err
	}

	return pool, nil
}
