package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"role-leader/internal/api"
)

//	type DB interface {
//		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
//	}
type Service struct {
	api.RoleLeaderServer
	logger *zap.Logger
	conn   *pgxpool.Pool
}

func New(logger *zap.Logger, conn *pgxpool.Pool) *Service {
	return &Service{
		logger: logger,
		conn:   conn,
	}
}

var (
	ErrEmptyMessage   = status.Errorf(codes.InvalidArgument, "Message cannot be empty")
	ErrCallIdNotFound = status.Errorf(codes.NotFound, "Call id not found")
	ErrInternalError  = status.Errorf(codes.Internal, "Internal server error")
)

func (s *Service) CreateFeedback(ctx context.Context, req *api.CreateFeedbackRequest) (*api.CreateFeedbackResponse, error) {
	if req.Message == "" {
		s.logger.Error("CreateFeedback: feedback message is empty")
		return nil,
			ErrEmptyMessage
	}

	q := "update schema_call.phone_call set feedback = $1 where call_id = $2"
	tag, err := s.conn.Exec(ctx, q, req.Message, req.CallId)

	if tag.RowsAffected() == 0 {
		s.logger.Error("CreateFeedback: call id not found: ", zap.String("call_id = ", req.CallId))
		return nil,
			ErrCallIdNotFound
	}
	if err != nil {
		s.logger.Error("CreateFeedback: failed to create feedback", zap.Error(err))
		return nil,
			ErrInternalError
	}

	s.logger.Info("Successfully created feedback", zap.String("feedback", req.Message))

	return &api.CreateFeedbackResponse{}, nil
}

func (s *Service) GetCall(ctx context.Context, req *api.GetCallRequest) (*api.GetCallResponse, error) {
	q := "select * from schema_call.phone_call where call_id = $1"
	var call api.Call

	err := s.conn.QueryRow(ctx, q, req.CallId).Scan(
		&call.CallId,
		&call.UserId,
		&call.LeaderId,
		&call.Title,
		&call.StartTime,
		&call.Status,
		&call.Feedback,
	)

	if err != nil {
		s.logger.Error("Failed to get call", zap.Error(err))
		return nil, fmt.Errorf("failed to get call: %w", err)
	}

	return &api.GetCallResponse{Call: &call}, nil
}

func (s *Service) GetLeaderCalls(ctx context.Context, req *api.GetLeaderCallsRequest) (*api.GetLeaderCallsResponse, error) {
	q := "select * from schema_call.phone_call where leader_id = $1"

	row, err := s.conn.Query(ctx, q, req.LeaderId)
	if err != nil {
		s.logger.Error("Failed to get calls", zap.Error(err))
		return nil, fmt.Errorf("failed to get calls: %w", err)
	}
	defer row.Close()

	var calls []*api.Call
	for row.Next() {
		var call api.Call
		err := row.Scan(
			&call.CallId,
			&call.UserId,
			&call.LeaderId,
			&call.Title,
			&call.StartTime,
			&call.Status,
			&call.Feedback,
		)
		if err != nil {
			s.logger.Error("Failed to get calls", zap.Error(err))
			return nil, fmt.Errorf("failed to get calls: %w", err)
		}
		calls = append(calls, &call)
	}

	if err := row.Err(); err != nil {
		s.logger.Error("Error after iterating calls", zap.Error(err))
		return nil, fmt.Errorf("error after iterating calls: %w", err)
	}
	for _, call := range calls {
		fmt.Println(call)
	}

	return &api.GetLeaderCallsResponse{Calls: calls}, nil
}
