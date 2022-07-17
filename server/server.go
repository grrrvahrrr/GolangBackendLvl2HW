package server

import (
	pb "GoBeLvl2/grpcauth"
	"GoBeLvl2/redis"
	"context"
	"net"

	"log"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthenticationServiceServer
	Redis *redis.RedisStore
}

func NewServer(redis *redis.RedisStore) *Server {
	return &Server{
		Redis: redis,
	}
}

func (s *Server) Register(ctx context.Context, in *pb.RegistrationRequest) (*pb.RegistrationReply, error) {
	login := in.GetLogin()
	password := in.GetPassword()

	authCode, err := s.Redis.CacheRegister(ctx, login, password)
	if err != nil {
		log.Println("redis error registering user", err)

	}
	return &pb.RegistrationReply{Authcode: authCode}, nil
}

func (s *Server) Authenticate(ctx context.Context, in *pb.AuthRequest) (*pb.AuthReply, error) {
	authCode := in.GetAuthcode()
	err := s.Redis.CacheAuth(ctx, authCode)
	if err != nil {
		log.Println("redis error authenticating user", err)
		return &pb.AuthReply{Response: "Authentication Failed!"}, err
	}

	return &pb.AuthReply{Response: "Authentication successful!"}, nil
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serv := grpc.NewServer()
	pb.RegisterAuthenticationServiceServer(serv, s)
	log.Printf("server listening at %v", lis.Addr())
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
