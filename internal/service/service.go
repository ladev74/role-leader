package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"role-leader/internal/api"
	"role-leader/internal/config"
)

type Service struct {
	api.RoleLeaderServer
	cfg    *config.Config
	logger *zap.Logger
	conn   *pgxpool.Pool
}

func New(cfg *config.Config, logger *zap.Logger, conn *pgxpool.Pool) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger,
		conn:   conn,
	}
}

func (s *Service) CreateFeedback(ctx context.Context, req *api.CreateFeedbackRequest) (*api.CreateFeedbackResponse, error) {
	q := "update schema_call.phone_call set feedback = $1 where call_id = $2"
	_, err := s.conn.Query(ctx, q, req.Message, req.CallId)

	if err != nil {
		s.logger.Error("CreateFeedback: failed to create feedback", zap.Error(err))
		return &api.CreateFeedbackResponse{Status: "ERROR"},
			fmt.Errorf("createFeedback: failed to create feedback: %w", err)
	}

	s.logger.Info("Successfully created feedback", zap.String("feedback", req.Message))

	return &api.CreateFeedbackResponse{Status: "OK"}, nil
}

func (s *Service) GetCall(ctx context.Context, req *api.GetCallRequest) (*api.GetCallResponse, error) {
	q := "select * from schema_call.phone_call where call_id = $1"
	t := time.Time{}
	var call api.Call

	err := s.conn.QueryRow(ctx, q, req.CallId).Scan(
		&call.CallId,
		&call.UserId,
		&call.LeaderId,
		&call.Title,
		&call.Status,
		&call.Feedback,
		&t,
	)
	call.StartTime = timestamppb.New(t)

	if err != nil {
		s.logger.Error("Failed to get call", zap.Error(err))
		return nil, fmt.Errorf("failed to get call: %w", err)
	}

	return &api.GetCallResponse{Call: &call}, nil
}

func (s *Service) ListCalls(ctx context.Context, req *api.ListCallsRequest) (*api.ListCallsResponse, error) {
	calls := make([]*api.Call, 0)
	return &api.ListCallsResponse{
		Calls: calls,
	}, nil
}
