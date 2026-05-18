package main

import (
	"context"
	"log"
	"net"

	pb "streaming-platform/proto/playback"
	grpcHandler "streaming-platform/services/playback-service/internal/delivery/grpc"
	postgresRepo "streaming-platform/services/playback-service/internal/repository/postgres"
	"streaming-platform/services/playback-service/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5435/playback_db?sslmode=disable",
	)
	if err != nil {
		log.Fatal("failed to connect to playback db:", err)
	}
	defer db.Close()

	repo := postgresRepo.NewPlaybackRepository(db)
	playbackUsecase := usecase.NewPlaybackUsecase(repo)
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