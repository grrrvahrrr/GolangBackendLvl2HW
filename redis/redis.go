package redis

import (
	"GoBeLvl2/entities"
	"GoBeLvl2/logic"
	"context"
	"encoding/json"
	"fmt"
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

	var value string
	err = rs.cache.Get(ctx, string(u), &value)
	if value != "registration_complete" {
		err = rs.cache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   authCode,
			Value: u,
			TTL:   time.Hour,
		})
		if err != nil {
			return "", err
		}
	} else {
		return "user_already_registred", err
	}

	return authCode, nil
}

func (rs *RedisStore) CacheAuth(ctx context.Context, authCode string) error {
	var value []byte
	exists := rs.cache.Exists(ctx, authCode)
	if !exists {
		return fmt.Errorf("authentication code not found, please, register again at /register")
	}
	err := rs.cache.Get(ctx, authCode, &value)
	if err != nil {
		return err
	}
	if err == nil && value != nil {
		err := rs.cache.Delete(ctx, authCode)
		if err != nil {
			return err
		}

		//Set user to cache as registered, but it needs to go to db of registered users
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
