package cacheport

import (
	"GoBeLvl2/recieveMsg/internal/entities"
	"context"

	"go.opentelemetry.io/otel/trace"
)

type CacheStore interface {
	CacheSet(ctx context.Context, key string, order entities.Order) error
	CacheGet(ctx context.Context, key string) (*entities.Order, error)
	CacheRestore(ctx context.Context, key string, orderId string) error
}

type CachePort struct {
	cs CacheStore
	Tr trace.Tracer
}

func NewCacheStorage(cs CacheStore, tr trace.Tracer) *CachePort {
	return &CachePort{
		cs: cs,
		Tr: tr,
	}
}

func (cp *CachePort) CacheSet(ctx context.Context, key string, order entities.Order) error {
	_, span := cp.Tr.Start(ctx, "RedisCacheSet")
	defer span.End()

	err := cp.cs.CacheSet(ctx, key, order)
	if err != nil {
		return err
	}
	return nil
}

func (cp *CachePort) CacheGet(ctx context.Context, key string) (*entities.Order, error) {
	_, span := cp.Tr.Start(ctx, "RedisCacheGet")
	defer span.End()

	order, err := cp.cs.CacheGet(ctx, key)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (cp *CachePort) CacheRestore(ctx context.Context, key string, orderId string) error {
	_, span := cp.Tr.Start(ctx, "RedisCacheRestore")
	defer span.End()

	err := cp.cs.CacheRestore(ctx, key, orderId)
	if err != nil {
		return err
	}
	return nil
}
