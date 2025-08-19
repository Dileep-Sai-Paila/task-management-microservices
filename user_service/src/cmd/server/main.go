package main

import (
	"log"
	"net"
	"net/http"
	"user_service/internal/adaptors/persistance"
	"user_service/internal/config"
	"user_service/internal/interfaces/input/api/rest/handler"
	"user_service/internal/interfaces/input/api/rest/middleware"
	grpcServer "user_service/internal/interfaces/output/grpc/server"
	"user_service/internal/usecase"
	pb "user_service/proto"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	dbStore, err := persistance.NewDBStore(cfg.DBSource)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	userUsecase := usecase.NewUserUsecase(dbStore, cfg.JWTSecretKey)

	// start grpc server in another goroutine
	go startGRPCServer(userUsecase, cfg.GRPCServerAddress)

	// REST server start
	userHandler := handler.NewUserHandler(userUsecase)
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Post("/register", userHandler.RegisterUser)
	r.Post("/login", userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecretKey))
		r.Get("/profile", userHandler.GetProfile)
	})

	log.Printf("REST Service starting on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}

func startGRPCServer(userUsecase usecase.UserUsecase, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	userServer := grpcServer.NewUserServer(userUsecase)
	pb.RegisterUserServiceServer(s, userServer)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
