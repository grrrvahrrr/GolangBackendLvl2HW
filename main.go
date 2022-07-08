package main

import (
	"GoBeLvl2/apigin"
	"GoBeLvl2/redis"
	"GoBeLvl2/server"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	//Creating Context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//Creating Cache Storage
	redis := redis.NewRedis()

	gr := apigin.NewApiGin(redis)

	srv := server.NewServer(":4444", gr)

	//Start Server
	srv.Start()

	//Hello
	fmt.Println("Recieving messages!")

	//Shutdown
	<-ctx.Done()
	srv.Stop()
	cancel()
}
