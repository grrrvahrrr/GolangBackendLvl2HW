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
	"GoBeLvl2/recieveMsg/internal/telemetry"
)

func main() {
	//Creating Context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//Tracing
	tp, err := telemetry.RunTracingCollection(ctx)
	if err != nil {
		log.Println("Error running tracing collection: ", err)
	}
	defer func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			log.Println("failed to stop the traces collector: ", err)
		}
	}()

	tr := tp.Tracer("server")

	//Creating DB Storage
	const dsn = "postgres://wbuser:wb@localhost:5433/wbl0db?sslmode=disable"
	udf, err := database.NewPgStorage(dsn)
	if err != nil {
		log.Println("Error creating database files: ", err)
	}

	db := dbport.NewDataStorage(udf, tr)

	//Creating Cache Storage
	redisStore := cachestore.NewRedis(db)
	cacheStore := cacheport.NewCacheStorage(redisStore, tr)

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
	router := api.NewApiChiRouter(h, tr)
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
