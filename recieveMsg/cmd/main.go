package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"GoBeLvl2/recieveMsg/internal/api"
	"GoBeLvl2/recieveMsg/internal/cacheport"
	"GoBeLvl2/recieveMsg/internal/cachestore"
	"GoBeLvl2/recieveMsg/internal/database"
	"GoBeLvl2/recieveMsg/internal/dbport"
	"GoBeLvl2/recieveMsg/internal/entities"
	"GoBeLvl2/recieveMsg/internal/natsstreaming"
	"GoBeLvl2/recieveMsg/internal/server"
	"GoBeLvl2/recieveMsg/internal/shardmanager"
)

func main() {
	//Creating Context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//Creating DB Storage
	m := shardmanager.NewManager()
	m.Add(&shardmanager.Shard{
		Address: "postgres://wbuser:wb@localhost:5433/wbl0db?sslmode=disable",
		Number:  0,
	})
	m.Add(&shardmanager.Shard{
		Address: "postgres://wbuser:wb@localhost:8110/wbl0db?sslmode=disable",
		Number:  1,
	})
	m.Add(&shardmanager.Shard{
		Address: "postgres://wbuser:wb@localhost:8120/wbl0db?sslmode=disable",
		Number:  2,
	})
	m.AddReplica(&shardmanager.Shard{
		Address: "postgres://wbuser:wb@localhost:8101/wbl0db?sslmode=disable",
		Number:  0,
	})
	m.AddReplica(&shardmanager.Shard{
		Address: "postgres://wbuser:wb@localhost:8111/wbl0db?sslmode=disable",
		Number:  1,
	})
	m.AddReplica(&shardmanager.Shard{
		Address: "postgres://wbuser:wb@localhost:8121/wbl0db?sslmode=disable",
		Number:  2,
	})

	udf, err := database.NewPgStorage(m)
	if err != nil {
		log.Println("Error creating database files: ", err)
	}

	db := dbport.NewDataStorage(udf)

	//Creating Cache Storage
	redisStore := cachestore.NewRedis(db)
	cacheStore := cacheport.NewCacheStorage(redisStore)

	//Nats streaming recieving msg
	ns, err := natsstreaming.NewNatsStreaming(db, cacheStore)
	if err != nil {
		log.Println("Error connecting to nats: ", err)
	}

	orderch := make(chan entities.Order)

	err = ns.ListenToNats(ctx, orderch)
	if err != nil {
		log.Println("Error listening to nats: ", err)
	}

	go ns.WritingMsgRoutine(ctx, orderch)

	//Front
	h := api.NewHandlers(cacheStore)
	router := api.NewApiChiRouter(h)
	srv := server.NewServer(":3333", router)

	//Start Server
	srv.Start()

	//Hello
	fmt.Println("Recieving messages!")

	//Shutdown
	<-ctx.Done()
	ns.SubClose()
	srv.Stop()
	cancel()
}
