package database

import (
	"GoBeLvl2/recieveMsg/internal/dbport"
	"GoBeLvl2/recieveMsg/internal/entities"
	"GoBeLvl2/recieveMsg/internal/shardmanager"
	"context"
	"database/sql"
	"encoding/json"
	"sync"

	_ "github.com/jackc/pgx/v4/stdlib" //psql driver
)

//Db data transition object can be added

var _ dbport.DbStore = &PgStorage{}

type PgStorage struct {
	db *sql.DB
	m  *shardmanager.Manager
	p  *Pool
}

//Connection pool
type Pool struct {
	sync.RWMutex
	cc map[string]*sql.DB
}

func NewPool() *Pool {
	return &Pool{
		cc: map[string]*sql.DB{},
	}
}
func (p *Pool) PoolConnection(addr string) (*sql.DB, error) {
	p.RLock()
	if c, ok := p.cc[addr]; ok {
		defer p.RUnlock()
		return c, nil
	}
	p.RUnlock()
	p.Lock()
	defer p.Unlock()
	if c, ok := p.cc[addr]; ok {
		return c, nil
	}
	var err error
	p.cc[addr], err = sql.Open("pgx", addr)
	return p.cc[addr], err
}

func NewPgStorage(m *shardmanager.Manager) (*PgStorage, error) {
	p := NewPool()

	us := &PgStorage{
		m: m,
		p: p,
	}

	return us, nil
}

func (pg *PgStorage) connection() (*sql.DB, error) {
	s, err := pg.m.Shard()
	if err != nil {
		return nil, err
	}
	return pg.p.PoolConnection(s.Address)
}

func (pg *PgStorage) GetOrderInfo(ctx context.Context, orderId string) (*entities.Order, error) {
	db, err := pg.connection()
	if err != nil {
		return nil, err
	}
	pg.db = db

	order := &entities.Order{
		OrderId: orderId,
	}

	var delivery []uint8
	var payment []uint8
	var items []uint8

	rows, err := pg.db.QueryContext(ctx, `SELECT track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM ordertest.orders WHERE order_uid = $1`, order.OrderId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(
			&order.TrackNumber, &order.Entry,
			&delivery, &payment, &items, &order.Locale,
			&order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.ShardKey,
			&order.SmID, &order.DateCreated, &order.OofShard,
		); err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}
	rows.Close()

	err = json.Unmarshal(delivery, &order.Delivery)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(payment, &order.Payment)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(items, &order.Items)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (pg *PgStorage) WriteOrder(ctx context.Context, order entities.Order) error {
	_, err := pg.db.ExecContext(ctx, `INSERT INTO ordertest.orders
	 (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderId, order.TrackNumber, order.Entry,
		order.Delivery, order.Payment, order.Items, order.Locale,
		order.InternalSignature, order.CustomerId, order.DeliveryService, order.ShardKey,
		order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStorage) WriteDelivery(ctx context.Context, order entities.Order) error {
	_, err := pg.db.ExecContext(ctx, `INSERT INTO ordertest.deliveries (order_uid, name, phone, zip, city, address, region, email)
	 VALUES($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderId, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStorage) WritePayment(ctx context.Context, order entities.Order) error {
	_, err := pg.db.ExecContext(ctx, `INSERT INTO ordertest.payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (transaction) DO NOTHING`,
		order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStorage) WriteCart(ctx context.Context, order entities.Order) error {
	_, err := pg.db.ExecContext(ctx, `INSERT INTO ordertest.carts (order_uid, items) VALUES ($1, $2) ON CONFLICT (order_uid) DO NOTHING`, order.OrderId, order.Items)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStorage) WriteOrderData(ctx context.Context, order entities.Order) error {
	db, err := pg.connection()
	if err != nil {
		return err
	}
	pg.db = db

	err = pg.WriteOrder(ctx, order)
	if err != nil {
		return err
	}

	err = pg.WriteDelivery(ctx, order)
	if err != nil {
		return err
	}

	err = pg.WritePayment(ctx, order)
	if err != nil {
		return err
	}

	err = pg.WriteCart(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
