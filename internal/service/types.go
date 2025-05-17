package service

import (
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
