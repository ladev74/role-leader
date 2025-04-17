package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Env string `yaml:"ENV" env:"ENV" envDefault:"local"`
}

func New(cfg Config) (*zap.Logger, error) {
	switch cfg.Env {
	case "local":
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("can't initialize logger: %v", err)
		}

		return logger, nil

	default:
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("can't initialize logger: %v", err)
		}

		return logger, nil
	}
}

func Interceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp any, err error) {

		logger.Info(
			"new request", zap.String("method", info.FullMethod),
			zap.Any("request", req),
			zap.Time("time", time.Now()),
		)

		return next(ctx, req)
	}
}
