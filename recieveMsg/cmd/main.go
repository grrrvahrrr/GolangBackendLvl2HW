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
	"GoBeLvl2/recieveMsg/internal/msgprocess"
	"GoBeLvl2/recieveMsg/internal/rabbit"
	"GoBeLvl2/recieveMsg/internal/server"
)

func main() {
	//Creating Context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//Creating DB Storage
	const dsn = "postgres://wbuser:wb@localhost:5433/wbl0db?sslmode=disable"
	udf, err := database.NewPgStorage(dsn)
	if err != nil {
		log.Println("Error creating database files: ", err)
	}

	db := dbport.NewDataStorage(udf)

	//Creating Cache Storage
	redisStore := cachestore.NewRedis(db)
	cacheStore := cacheport.NewCacheStorage(redisStore)

	//MsgRecieving
	orderch := make(chan entities.Order)

	//Nats streaming recieving msg
	// ns, err := natsstreaming.NewNatsStreaming()
	// if err != nil {
	// 	log.Println("Error connecting to nats: ", err)
	// }

	// err = ns.ListenToNats(ctx, orderch)
	// if err != nil {
	// 	log.Println("Error listening to nats: ", err)
	// }

	//Rabbit
	r, err := rabbit.NewRabbit(ctx)
	if err != nil {
		log.Println("Error connecting to Rabbit: ", err)
	}

	go r.ListenToRabbit(orderch)

	//Kafka
	// k, err := kafka.NewKafka(ctx)
	// if err != nil {
	// 	log.Println("Error connecting to Kafka: ", err)
	// }

	// go k.ListenToKafka(orderch)

	//Message processing goroutine
	mp := msgprocess.NewMsgProcess(db, cacheStore)
	go mp.WritingMsgRoutine(ctx, orderch)

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
	//ns.SubClose()
	r.SubClose()
	//k.SubClose()
	srv.Stop()
	cancel()
}
