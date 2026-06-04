package database

import (
	"errors"
	"strconv"

	"github.com/RaflyAdiyasa/Helpdesk-Ticketing-API/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Config) (*redis.Client, error) {
	DB, _ := strconv.Atoi(cfg.Redis.DB)
	client := redis.NewClient(&redis.Options{
		DB:   DB,
		Addr: cfg.Redis.Endpoint,
	})

	if client != nil {
		return client, nil
	} else {
		return nil, errors.New("Tidak bisa connect Redis")
	}

}
