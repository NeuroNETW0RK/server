package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

type Options struct {
	Addr       string
	Port       int
	Db         int
	MaxReTries int
	Password   string
	PoolSize   int
}

func New(opts *Options) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%v", opts.Addr, opts.Port),
		Password:   opts.Password,
		DB:         opts.Db,
		MaxRetries: opts.MaxReTries,
		PoolSize:   opts.PoolSize,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return client, nil
}
