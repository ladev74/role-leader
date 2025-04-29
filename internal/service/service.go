package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"role-leader/internal/api"
)

var (
	ErrEmptyMessage     = status.Errorf(codes.InvalidArgument, "Message cannot be empty")
	ErrCallIdNotFound   = status.Errorf(codes.NotFound, "Call id not found")
	ErrLeaderIdNotFound = status.Errorf(codes.NotFound, "Leader id not found")
	ErrInternalError    = status.Errorf(codes.Internal, "Internal server error")
)

const (
	QueryForCreateFeedback = "update schema_call.phone_call set feedback = $1 where call_id = $2"
	QueryForGetCall        = "select * from schema_call.phone_call where call_id = $1"
	QueryForGetLeaderCalls = "select * from schema_call.phone_call where leader_id = $1"
)

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

func (s *Service) CreateFeedback(ctx context.Context, req *api.CreateFeedbackRequest) (*api.CreateFeedbackResponse, error) {
	if req.Message == "" {
		s.logger.Error("CreateFeedback: feedback message is empty", zap.String("call_id = ", req.CallId))
		return nil,
			ErrEmptyMessage
	}

	tag, err := s.conn.Exec(ctx, QueryForCreateFeedback, req.Message, req.CallId)

	if err != nil {
		s.logger.Error("CreateFeedback: failed to create feedback", zap.Error(err), zap.String("call_id = ", req.CallId))
		return nil,
			ErrInternalError
	}

	if tag.RowsAffected() == 0 {
		s.logger.Error("CreateFeedback: call id not found: ", zap.String("call_id = ", req.CallId))
		return nil,
			ErrCallIdNotFound
	}

	s.logger.Info("Successfully created feedback", zap.String("call_id = ", req.CallId), zap.String("feedback", req.Message))

	return &api.CreateFeedbackResponse{}, nil
}

func (s *Service) GetCall(ctx context.Context, req *api.GetCallRequest) (*api.GetCallResponse, error) {
	var call api.Call

	err := s.conn.QueryRow(ctx, QueryForGetCall, req.CallId).Scan(
		&call.CallId,
		&call.UserId,
		&call.LeaderId,
		&call.Title,
		&call.Status,
		&call.Feedback,
		&call.StartTime,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("GetCall: call id not found: ", zap.Error(err), zap.String("call_id = ", req.CallId))
		return nil, ErrCallIdNotFound
	}
	if err != nil {
		s.logger.Error("GetCall: failed to get call", zap.Error(err))
		return nil, ErrInternalError
	}

	return &api.GetCallResponse{Call: &call}, nil
}

func (s *Service) GetLeaderCalls(ctx context.Context, req *api.GetLeaderCallsRequest) (*api.GetLeaderCallsResponse, error) {
	row, err := s.conn.Query(ctx, QueryForGetLeaderCalls, req.LeaderId)

	if err != nil {
		s.logger.Error("GetLeaderCalls: failed to get calls", zap.Error(err), zap.String("leader id = ", req.LeaderId))
		return nil, ErrInternalError
	}

	defer row.Close()

	if errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("GetLeaderCalls: leader id not found: ", zap.String("leader id = ", req.LeaderId))
		return nil, ErrLeaderIdNotFound
	}

	var calls []*api.Call
	for row.Next() {
		var call api.Call
		err := row.Scan(
			&call.CallId,
			&call.UserId,
			&call.LeaderId,
			&call.Title,
			&call.Status,
			&call.Feedback,
			&call.StartTime,
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

	if len(calls) == 0 {
		s.logger.Error("GetLeaderCalls: leader id no found", zap.String("leader id = ", req.LeaderId))
		return nil, ErrLeaderIdNotFound
	}

	return &api.GetLeaderCallsResponse{Calls: calls}, nil
}
