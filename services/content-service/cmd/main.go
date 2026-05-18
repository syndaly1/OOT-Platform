package main

import (
	"database/sql"
	"log"
	"net"

	"content-service/internal/config"
	contentgrpc "content-service/internal/delivery/grpc"
	"content-service/internal/events"
	"content-service/internal/repository/postgres"
	"content-service/internal/usecase"
	"content-service/pkg/contentpb"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	log.Println("connected to PostgreSQL")

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	log.Println("migrations applied")

	publisher, err := events.NewPublisher(cfg.NATSURL)
	if err != nil {
		log.Printf("failed to connect to NATS: %v", err)
	} else {
		defer publisher.Close()
		log.Println("connected to NATS")
	}

	contentRepo := postgres.NewContentPostgresRepository(db)
	episodeRepo := postgres.NewEpisodePostgresRepository(db)

	contentUsecase := usecase.NewContentUsecase(contentRepo, publisher)
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
