package middleware

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/emitra-labs/common/errors"
	"github.com/emitra-labs/common/log"
	"github.com/emitra-labs/gorest/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

type RateLimiterRedisStore struct {
	Count  int
	Period time.Duration
}

func (r RateLimiterRedisStore) Allow(identifier string) (bool, error) {
	if store.RedisClient == nil || store.RedisLock == nil {
		log.Error("Looks like you forgot to set GOREST_REDIS_URL")
		return false, errors.Internal()
	}

	ctx := context.Background()

	lock, err := store.RedisLock.Obtain(ctx, "lock_"+identifier, 100*time.Millisecond, nil)

	if err != nil {
		if err != redislock.ErrNotObtained {
			log.Errorf("Failed to obtain lock: %s", err)
		}

		return false, nil
	}

	defer lock.Release(ctx)

	val, err := store.RedisClient.Get(ctx, identifier).Int()
	if err != nil {
		if err == redis.Nil {
			store.RedisClient.Set(ctx, identifier, 1, r.Period)
			return true, nil
		}

		log.Errorf("Failed to get value from Redis: %s", err)
		return false, nil
	}

	if val >= r.Count {
		return false, nil
	} else {
		store.RedisClient.Incr(ctx, identifier)
	}

	return true, nil
}

func RateLimit(count int, period time.Duration) echo.MiddlewareFunc {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: &RateLimiterRedisStore{
			Count:  count,
			Period: period,
		},
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return errors.PermissionDenied("Too many requests")
		},
	}

	return middleware.RateLimiterWithConfig(config)
}
