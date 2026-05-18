package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "streaming-platform/proto/playback"
	grpcHandler "streaming-platform/services/playback-service/internal/delivery/grpc"
	natsPublisher "streaming-platform/services/playback-service/internal/event/nats"
	postgresRepo "streaming-platform/services/playback-service/internal/repository/postgres"
	"streaming-platform/services/playback-service/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	natsgo "github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, playbackDatabaseURL())
	if err != nil {
		log.Fatal("failed to connect to playback db:", err)
	}
	defer db.Close()

	natsConn, err := natsgo.Connect(getenv("NATS_URL", "nats://localhost:4222"))
	if err != nil {
		log.Fatal("failed to connect to nats:", err)
	}
	defer natsConn.Close()

	repo := postgresRepo.NewPlaybackRepository(db)
	publisher := natsPublisher.NewPublisher(natsConn)
	playbackUsecase := usecase.NewPlaybackUsecase(repo, publisher)
	handler := grpcHandler.NewPlaybackHandler(playbackUsecase)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	pb.RegisterPlaybackServiceServer(server, handler)
	reflection.Register(server)

	log.Println("Playback Service running on :50053")

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func playbackDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getenv("PLAYBACK_DB_USER", "postgres"),
		getenv("PLAYBACK_DB_PASSWORD", "postgres"),
		getenv("PLAYBACK_DB_HOST", "localhost"),
		getenv("PLAYBACK_DB_PORT", "5432"),
		getenv("PLAYBACK_DB_NAME", "playback_db"),
	)
}
