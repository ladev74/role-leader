package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"role-leader/internal/api"
	"role-leader/internal/service"
)

// добавить реализацию с sqlite
func TestCreateFeedback(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal("error creating mock pool:", err)
	}
	defer mock.Close()

	srv := service.New(zap.NewNop(), mock)

	tests := []struct {
		name       string
		req        api.CreateFeedbackRequest
		mockSetup  func()
		wantStatus string
		wantErr    bool
	}{
		{
			name: "successful update",
			req: api.CreateFeedbackRequest{
				CallId:  "1111",
				Message: "aboba",
			},
			mockSetup: func() {
				mock.ExpectExec("update schema_call.phone_call").
					WithArgs("aboba", "1111").
					WillReturnResult(pgxmock.NewResult("update", 1))
			},
			wantStatus: service.StatusOkForGrpcResponse,
			wantErr:    false,
		},
		{
			name: "empty message",
			req: api.CreateFeedbackRequest{
				CallId:  "1111",
				Message: "",
			},
			mockSetup:  func() {},
			wantStatus: service.StatusErrForGrpcResponse,
			wantErr:    true,
		},
		{
			name: "database error",
			req: api.CreateFeedbackRequest{
				CallId:  "3333",
				Message: "test message 3",
			},
			mockSetup: func() {
				mock.ExpectExec("update schema_call.phone_call set feedback = \\$1 where call_id = \\$2").
					WithArgs("test message 3", "3333").
					WillReturnError(fmt.Errorf("database error"))
			},
			wantStatus: service.StatusErrForGrpcResponse,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			got, err := srv.CreateFeedback(ctx, &tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStatus, got.Status)
			assert.NoError(t, mock.ExpectationsWereMet())

			mock.Reset()
		})
	}

}
