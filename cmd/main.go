package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"role-leader/internal/api"
	"role-leader/internal/config"
	"role-leader/internal/logger"
	"role-leader/internal/postgres"
	"role-leader/internal/service"
)

func main() {
	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't initialize config: %v", err)
	}

	l, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("can't initialize l: %v", err)
	}

	conn, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		l.Fatal("can't initialize database", zap.Error(err))
	}
	err = conn.Ping(ctx)
	if err != nil {
		l.Fatal("can't initialize database", zap.Error(err))
	}
	l.Info("connected to postgres")

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPC.Port))
	if err != nil {
		l.Fatal("failed to listen", zap.Error(err))
	}

	srv := service.New(l, conn)

	server := grpc.NewServer(grpc.UnaryInterceptor(logger.Interceptor(l)))
	api.RegisterRoleLeaderServer(server, srv)
	reflection.Register(server)

	go func() {
		if err := server.Serve(lis); err != nil {
			l.Fatal("failed to serve", zap.Error(err))
		}
	}()

	l.Info("grpc server started")

	select {
	case <-ctx.Done():
		server.GracefulStop()
		conn.Close()
		l.Info("grpc server stopped")
	}

}
