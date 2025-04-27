package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"role-leader/internal/api"
)

type Service struct {
	api.RoleLeaderServer
	logger *zap.Logger
	db     DB
}

func New(logger *zap.Logger, db DB) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}

const (
	StatusOkForGrpcResponse  = "OK"
	StatusErrForGrpcResponse = "ERROR"
)

type DB interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (s *Service) CreateFeedback(ctx context.Context, req *api.CreateFeedbackRequest) (*api.CreateFeedbackResponse, error) {
	q := "update schema_call.phone_call set feedback = $1 where call_id = $2"
	_, err := s.db.Exec(ctx, q, req.Message, req.CallId)

	if err != nil {
		s.logger.Error("CreateFeedback: failed to create feedback", zap.Error(err))
		return &api.CreateFeedbackResponse{Status: StatusErrForGrpcResponse},
			fmt.Errorf("createFeedback: failed to create feedback: %w", err)
	}

	s.logger.Info("Successfully created feedback", zap.String("feedback", req.Message))

	return &api.CreateFeedbackResponse{Status: StatusOkForGrpcResponse}, nil
}

func (s *Service) GetCall(ctx context.Context, req *api.GetCallRequest) (*api.GetCallResponse, error) {
	q := "select * from schema_call.phone_call where call_id = $1"
	t := time.Time{}
	var call api.Call

	err := s.db.QueryRow(ctx, q, req.CallId).Scan(
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

func (s *Service) GetLeaderCalls(ctx context.Context, req *api.GetLeaderCallsRequest) (*api.GetLeaderCallsResponse, error) {
	q := "select * from schema_call.phone_call where leader_id = $1"

	row, err := s.db.Query(ctx, q, req.LeaderId)
	if err != nil {
		s.logger.Error("Failed to get calls", zap.Error(err))
		return nil, fmt.Errorf("failed to get calls: %w", err)
	}
	defer row.Close()

	var calls []*api.Call
	for row.Next() {
		var call api.Call
		var t time.Time
		err := row.Scan(
			&call.CallId,
			&call.UserId,
			&call.LeaderId,
			&call.Title,
			&call.Status,
			&call.Feedback,
			&t,
		)
		if err != nil {
			s.logger.Error("Failed to get calls", zap.Error(err))
			return nil, fmt.Errorf("failed to get calls: %w", err)
		}
		call.StartTime = timestamppb.New(t)
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
