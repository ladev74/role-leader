package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"role-leader/internal/api"
	"role-leader/internal/config"
)

func RunRest(ctx context.Context, l *zap.Logger, cfg config.Config) error {
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
	httpMux.Handle("/", LoggingMiddleware(l, mux))

	restServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.REST.Port),
		Handler: httpMux,
	}

	if err := restServer.ListenAndServe(); err != nil {
		l.Fatal("failed to serve REST", zap.Error(err))
		return fmt.Errorf("failed to serve REST: %w", err)
	}

	return nil
}

func LoggingMiddleware(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info(
			"new request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote-addr", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.Time("time", time.Now()),
		)

		w.Header().Set("X-Status", "OK")
		next.ServeHTTP(w, r)

	})
}
