package msgprocess

import (
	"GoBeLvl2/recieveMsg/internal/cacheport"
	"GoBeLvl2/recieveMsg/internal/dbport"
	"GoBeLvl2/recieveMsg/internal/entities"
	"context"
	"fmt"
	"log"
)

type MsgProcess struct {
	db         *dbport.DbPort
	cachestore *cacheport.CachePort
}

func NewMsgProcess(db *dbport.DbPort, cachestore *cacheport.CachePort) *MsgProcess {
	return &MsgProcess{
		db:         db,
		cachestore: cachestore,
	}
}

func (ns *MsgProcess) WriteMsgToDbAndCache(ctx context.Context, order entities.Order) error {
	//Writing order to DB
	err := ns.db.WriteOrderData(ctx, order)
	if err != nil {
		return err
	}

	//Writing to cache
	key := fmt.Sprintf("order:%s", order.OrderId)
	err = ns.cachestore.CacheSet(ctx, key, order)
	if err != nil {
		return err
	}
	return nil
}

func (ns *MsgProcess) WritingMsgRoutine(ctx context.Context, orderch chan entities.Order) {
	var order entities.Order
	for i := range orderch {
		order = i
		if order.OrderId != "" {
			err := ns.WriteMsgToDbAndCache(ctx, order)
			if err != nil {
				log.Fatal("Error writing to db or redis: ", err)
			}
		}
	}
}
