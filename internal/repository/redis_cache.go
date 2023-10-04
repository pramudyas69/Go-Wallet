package repository

import (
	"context"
	"e-wallet/domain"
	"e-wallet/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisRepository struct {
	rdb *redis.Client
}

func NewRedisClient(cnf *config.Config) domain.CacheRepository {
	return &redisRepository{
		rdb: redis.NewClient(&redis.Options{
			Addr:     cnf.Redis.Addr,
			Password: cnf.Redis.Pass,
			DB:       0,
		}),
	}
}

func (r redisRepository) Get(key string) ([]byte, error) {
	val, err := r.rdb.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r redisRepository) Set(key string, entry []byte) error {
	return r.rdb.Set(context.Background(), key, entry, 15*time.Minute).Err()
}

func (r redisRepository) Delete(key string) error {
	return r.rdb.Del(context.Background(), key).Err()
}
