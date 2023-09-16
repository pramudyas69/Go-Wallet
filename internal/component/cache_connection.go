package component

import (
	"context"
	"e-wallet/domain"
	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func GetCacheConnection() domain.CacheRepository {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Fatalf("error when connecting cache %s", err.Error())
	}
	return cache
}
