package redis

import (
	"GoBeLvl2/entities"
	"GoBeLvl2/logic"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	cache *cache.Cache
}

func NewRedis() *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &RedisStore{
		cache: rCache,
	}
}

func (rs *RedisStore) CacheRegister(ctx context.Context, login string, password string) (string, error) {
	var user entities.User
	user.Login = login
	user.Password = password
	u, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	authCode := logic.RandStringRunes(10)

	err = rs.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   authCode,
		Value: u,
		TTL:   time.Hour,
	})
	if err != nil {
		return "", err
	}

	return authCode, nil
}

func (rs *RedisStore) CacheAuth(ctx context.Context, authCode string) error {
	var value []byte
	err := rs.cache.Get(ctx, authCode, &value)
	if err != nil {
		return err
	}
	if err == nil && value != nil {
		err := rs.cache.Delete(ctx, authCode)
		if err != nil {
			return err
		}
		err = rs.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   string(value),
			Value: "registration_complete",
			TTL:   time.Hour,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
