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
	"role-leader/internal/service/gateway"
	"role-leader/internal/service/grpcSrv"
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
		log.Fatalf("can't initialize logger: %v", err)
	}
	defer l.Sync()

	conn, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		l.Fatal("can't initialize database", zap.Error(err))
	}

	if err = conn.Ping(ctx); err != nil {
		l.Fatal("can't initialize database", zap.Error(err))
	}
	l.Info("connected to postgres")

	grpcLis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPC.Port))
	if err != nil {
		l.Fatal("failed to listen", zap.Error(err))
	}
	srv := grpcSrv.New(l, conn)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(logger.Interceptor(l)))
	api.RegisterRoleLeaderServer(grpcServer, srv)
	reflection.Register(grpcServer)

	go func() {
		l.Info("GRPC server started", zap.Int("addr: ", cfg.GRPC.Port))
		if err := grpcServer.Serve(grpcLis); err != nil {
			l.Fatal("failed to serve", zap.Error(err))
		}
	}()

	go func() {
		l.Info("REST server started", zap.Int("addr: ", cfg.REST.Port))
		if err := gateway.RunRest(ctx, l, *cfg); err != nil {
			l.Fatal("failed to serve", zap.Error(err))
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		conn.Close()
		l.Info("GRPC server stopped", zap.Int("addr: ", cfg.GRPC.Port))
	}
}
