package main

import (
	"database/sql"
	"log"
	"net"

	"content-service/internal/config"
	contentgrpc "content-service/internal/delivery/grpc"
	"content-service/internal/repository/postgres"
	"content-service/internal/usecase"
	"content-service/pkg/contentpb"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	contentRepo := postgres.NewContentPostgresRepository(db)
	episodeRepo := postgres.NewEpisodePostgresRepository(db)

	contentUsecase := usecase.NewContentUsecase(contentRepo)
	episodeUsecase := usecase.NewEpisodeUsecase(episodeRepo)

	handler := contentgrpc.NewContentHandler(contentUsecase, episodeUsecase)

	server := grpc.NewServer()
	contentpb.RegisterContentServiceServer(server, handler)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Content Service is running on port %s", cfg.GRPCPort)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc server: %v", err)
	}
}
