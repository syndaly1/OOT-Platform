package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"streaming-platform/auth-service/config"
	grpchandler "streaming-platform/auth-service/internal/delivery/grpc"
	"streaming-platform/auth-service/internal/repository/postgres"
	rediscache "streaming-platform/auth-service/internal/repository/redis"
	"streaming-platform/auth-service/internal/usecase"
	pb "streaming-platform/proto/auth"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type userRegisteredEvent struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

type passwordResetEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func main() {
	cfg := config.Load()

	dbDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
	)

	migrateDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatal("db open:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("db ping:", err)
	}
	log.Println("connected to PostgreSQL")

	m, err := migrate.New("file://migrations", migrateDSN)
	if err != nil {
		log.Fatal("migrate init:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate up:", err)
	}
	log.Println("migrations applied")

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis connect:", err)
	}
	defer redisClient.Close()
	log.Println("connected to Redis")

	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatal("nats connect:", err)
	}
	defer nc.Close()
	log.Println("connected to NATS")

	accountRepo := postgres.NewAccountRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	emailVerRepo := postgres.NewEmailVerificationRepository(db)
	passwordResetRepo := postgres.NewPasswordResetRepository(db)

	profileCache := rediscache.NewProfileCache(redisClient)

	emailSender := usecase.NewEmailSender(cfg)

	nc.Subscribe("user.registered", func(msg *nats.Msg) {
		var event userRegisteredEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("user.registered parse error: %v", err)
			return
		}
		if err := emailSender.SendVerificationEmail(event.Email, event.Token); err != nil {
			log.Printf("send verification email error: %v", err)
			return
		}
		log.Printf("verification email sent to %s", event.Email)
	})

	nc.Subscribe("password.reset.requested", func(msg *nats.Msg) {
		var event passwordResetEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("password.reset.requested parse error: %v", err)
			return
		}
		if err := emailSender.SendPasswordResetEmail(event.Email, event.Token); err != nil {
			log.Printf("send password reset email error: %v", err)
			return
		}
		log.Printf("password reset email sent to %s", event.Email)
	})

	authUC := usecase.NewAuthUsecase(
		accountRepo,
		refreshTokenRepo,
		emailVerRepo,
		passwordResetRepo,
		nc,
		cfg,
	)
	accountUC := usecase.NewAccountUsecase(accountRepo)
	profileUC := usecase.NewProfileUsecase(profileRepo, profileCache)

	handler := grpchandler.NewAuthHandler(authUC, accountUC, profileUC)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("listen:", err)
	}

	log.Printf("auth-service gRPC listening on :%s", cfg.GRPC.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("serve:", err)
	}
}
