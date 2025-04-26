package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"role-leader/internal/api"
)

func TestService_CreateFeedback(t *testing.T) {
	type fields struct {
		RoleLeaderServer api.RoleLeaderServer
		cfg              *config.Config
		logger           *zap.Logger
		conn             *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		req *api.CreateFeedbackRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.CreateFeedbackResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				RoleLeaderServer: tt.fields.RoleLeaderServer,
				cfg:              tt.fields.cfg,
				logger:           tt.fields.logger,
				conn:             tt.fields.conn,
			}
			got, err := s.CreateFeedback(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFeedback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFeedback() got = %v, want %v", got, tt.want)
			}
		})
	}
}
