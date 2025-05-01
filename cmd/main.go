package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	srv := service.New(l, conn)

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
		if err := runRest(ctx, l, *cfg); err != nil {
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

func runRest(ctx context.Context, l *zap.Logger, cfg config.Config) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := api.RegisterRoleLeaderHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("0.0.0.0:%d", cfg.GRPC.Port),
		opts,
	)

	if err != nil {
		l.Fatal("failed to register REST handlers", zap.Error(err))
		return fmt.Errorf("failed to register REST handlers: %w", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)

	restServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.REST.Port),
		Handler: httpMux,
		//Handler: loggingMiddleware(l, httpMux),
	}

	if err := restServer.ListenAndServe(); err != nil {
		l.Fatal("failed to serve REST", zap.Error(err))
		return fmt.Errorf("failed to serve REST: %w", err)
	}

	return nil
}
func loggingMiddleware(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info(
			"new request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.Time("time", time.Now()),
		)
		//wrappedWriter := &responseLogger{w: w, status: http.StatusOK}

		next.ServeHTTP(w, r)
	})
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	return l.w.Write(b)
}

func (l *responseLogger) WriteHeader(statusCode int) {
	l.status = statusCode
	l.w.WriteHeader(statusCode)
}
