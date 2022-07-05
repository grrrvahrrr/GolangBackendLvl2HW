package dbport

import (
	"GoBeLvl2/recieveMsg/internal/entities"
	"context"

	"go.opentelemetry.io/otel/trace"
)

type DbStore interface {
	WriteOrderData(ctx context.Context, order entities.Order) error
	GetOrderInfo(ctx context.Context, orderId string) (*entities.Order, error)
}

type DbPort struct {
	dbstore DbStore
	tr      trace.Tracer
}

func NewDataStorage(dbstore DbStore, tr trace.Tracer) *DbPort {
	return &DbPort{
		dbstore: dbstore,
		tr:      tr,
	}
}

func (dp *DbPort) WriteOrderData(ctx context.Context, order entities.Order) error {
	_, span := dp.tr.Start(ctx, "DbWriteOrderData")
	defer span.End()

	err := dp.dbstore.WriteOrderData(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (dp *DbPort) GetOrderInfo(ctx context.Context, orderId string) (*entities.Order, error) {
	_, span := dp.tr.Start(ctx, "DbGetOrderInfo")
	defer span.End()
	order, err := dp.dbstore.GetOrderInfo(ctx, orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
