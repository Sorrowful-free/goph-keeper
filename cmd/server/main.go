package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gophkeeper/gophkeeper/internal/migrations"
	"github.com/gophkeeper/gophkeeper/internal/repository"
	"github.com/gophkeeper/gophkeeper/internal/server"
	"github.com/gophkeeper/gophkeeper/internal/storage"
	"github.com/gophkeeper/gophkeeper/internal/usecase/auth"
	"github.com/gophkeeper/gophkeeper/internal/usecase/data"
	"github.com/gophkeeper/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port     = flag.String("port", "50051", "Server port")
	dsn      = flag.String("dsn", "", "Database connection string (default: SQLite)")
	grpcAddr = flag.String("addr", "", "gRPC server address (overrides port)")
)

func main() {
	flag.Parse()

	addr := *grpcAddr
	if addr == "" {
		addr = fmt.Sprintf(":%s", *port)
	}

	// Infrastructure: storage (БД)
	st, err := storage.NewStorage(*dsn)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer st.Close()

	// Миграции (go-migrate)
	if err := migrations.RunUp(st.GetDB(), *dsn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Repositories (адаптеры к storage)
	userRepo := repository.NewUserRepository(st)
	dataRepo := repository.NewDataRepository(st)

	// Use cases
	authUC := auth.NewAuthUseCase(userRepo)
	dataUC := data.NewDataUseCase(dataRepo)

	// Delivery: gRPC services
	authService := server.NewAuthService(authUC)
	dataService := server.NewDataService(dataUC)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(server.LoggingInterceptor, server.AuthInterceptor),
	)

	proto.RegisterAuthServiceServer(grpcServer, authService)
	proto.RegisterDataServiceServer(grpcServer, dataService)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", addr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
