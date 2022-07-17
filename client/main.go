package main

import (
	"context"
	"log"
	"time"

	pb "GoBeLvl2/grpcauth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultLogin    = "qwe"
	defaultPassword = "qwe"
	authCode        = "iZltLLpAMu"
)

func main() {
	conn, err := grpc.Dial("localhost:3333", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthenticationServiceClient(conn)

	//register(c)
	auth(c)

}

func register(c pb.AuthenticationServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Register(ctx, &pb.RegistrationRequest{Login: defaultLogin, Password: defaultPassword})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("AuthCode: %s", r.GetAuthcode())
}

func auth(c pb.AuthenticationServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Authenticate(ctx, &pb.AuthRequest{Authcode: authCode})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf(r.GetResponse())
}
