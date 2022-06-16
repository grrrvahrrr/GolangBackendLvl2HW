package main

import (
	"GoBeLvl2/internal/api/handlers"
	"GoBeLvl2/internal/api/openapichi"
	"GoBeLvl2/internal/database/dbport"
	"GoBeLvl2/internal/database/psql"
	"GoBeLvl2/internal/processing"
	"GoBeLvl2/internal/server"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	//Precommit
	//Creating Context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	//Creating Storage
	const dsn = "postgres://user:123@localhost:5433/cataloguedb?sslmode=disable"
	udf, err := psql.NewPgStorage(dsn)
	if err != nil {
		log.Fatal("Error creating database files: ", err)
	}

	dbport := dbport.NewDbStorage(udf)

	//Processing
	p := processing.NewProcessing(dbport)
	//Handlers
	h := handlers.NewHandlers(p)
	//Router
	r := openapichi.NewOpenApiRouter(h)
	//Server
	srv := server.NewServer(":8000", r)

	//Starting
	srv.Start()

	//Shutting down
	<-ctx.Done()
	srv.Stop()

	fmt.Print("Server shutdown.")

}
