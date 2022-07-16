package main

import (
	"GoBeLvl2/apigin"
	"GoBeLvl2/elastic"
	"GoBeLvl2/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	//Creating Context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	//Elastic
	el, err := elastic.New([]string{"http://0.0.0.0:9200"})
	if err != nil {
		log.Fatalln(err)
	}
	if err := el.CreateIndex("post"); err != nil {
		log.Fatalln(err)
	}

	// Bootstrap storage.
	storage, err := elastic.NewPostStorage(*el)
	if err != nil {
		log.Fatalln(err)
	}

	//router and server
	gr := apigin.NewApiGin(storage)

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
