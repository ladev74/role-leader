package servicee

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"role-leader/internal/api"
	"role-leader/internal/config"
)

type Service struct {
	api.RoleLeaderServer
	cfg    *config.Config
	logger *zap.Logger
	conn   *pgx.Conn
}

func New(cfg *config.Config, logger *zap.Logger, conn *pgx.Conn) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger,
		conn:   conn,
	}
}

func (s *Service) CreateFeedback(ctx context.Context, request *api.CreateFeedbackRequest) (*api.CreateFeedbackResponse, error) {
	return nil, nil
}
