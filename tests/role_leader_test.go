package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"role-leader/internal/api"
	"role-leader/internal/postgres"
	"role-leader/internal/service"
)

var (
	connP *pgxpool.Pool
)

func upDB(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := postgres.Config{
		Username: "root",
		Password: "1234",
		Database: "postgres",
	}

	req := testcontainers.ContainerRequest{
		Name:         "postgres-for-tests",
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp", "8025/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": cfg.Password,
			"POSTGRES_USER":     cfg.Username,
			"POSTGRES_DB":       cfg.Database,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	cfg.Port, cfg.Host = port.Port(), host

	cfgForMigration := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	wd, _ := os.Getwd()

	migration, err := migrate.New("file://"+filepath.Join(wd, "../storage/migrations-for-tests"), cfgForMigration)
	if err != nil {
		return nil, err
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	cfgForPool := cfgForMigration + "&pool_max_conns=10&pool_min_conns=5"

	conn, err := pgxpool.New(ctx, cfgForPool)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	connP, err = upDB(ctx)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestCreateFeedback(t *testing.T) {
	ctx := context.Background()

	srv := service.New(zap.NewNop(), connP)

	tests := []struct {
		name    string
		req     *api.CreateFeedbackRequest
		want    *api.Call
		wantErr error
	}{
		{
			name: "successful create feedback",
			req: &api.CreateFeedbackRequest{
				CallId:  "1111",
				Message: "aboba",
			},
			want: &api.Call{
				CallId:    "1111",
				UserId:    "user1",
				LeaderId:  "leader1",
				Title:     "title1",
				Status:    "status1",
				Feedback:  "aboba",
				StartTime: "01:01:01",
			},
			wantErr: nil,
		},
		{
			name: "empty feedback",
			req: &api.CreateFeedbackRequest{
				CallId:  "2222",
				Message: "",
			},
			want: &api.Call{
				CallId:    "2222",
				UserId:    "user2",
				LeaderId:  "leader2",
				Title:     "title2",
				Status:    "status2",
				Feedback:  "feedback2",
				StartTime: "02:02:02",
			},
			wantErr: service.ErrEmptyMessage,
		},
		{
			name: "not existing call id",
			req: &api.CreateFeedbackRequest{
				CallId:  "0000",
				Message: "aboba",
			},
			want:    &api.Call{},
			wantErr: service.ErrCallIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "select * from schema_call.phone_call where call_id = $1"
			var got api.Call

			_, err := srv.CreateFeedback(ctx, tt.req)
			connP.QueryRow(ctx, q, tt.req.CallId).Scan(
				&got.CallId,
				&got.UserId,
				&got.LeaderId,
				&got.Title,
				&got.Status,
				&got.Feedback,
				&got.StartTime,
			)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, &got, tt.want)
		})
	}
}

func TestGetCall(t *testing.T) {
	ctx := context.Background()

	srv := service.New(zap.NewNop(), connP)

	tests := []struct {
		name    string
		req     *api.GetCallRequest
		want    *api.GetCallResponse
		wantErr error
	}{
		{
			name: "successful get call",
			req: &api.GetCallRequest{
				CallId: "2222",
			},

			want: &api.GetCallResponse{
				Call: &api.Call{
					CallId:    "2222",
					UserId:    "user2",
					LeaderId:  "leader2",
					Title:     "title2",
					Status:    "status2",
					Feedback:  "feedback2",
					StartTime: "02:02:02",
				},
			},
			wantErr: nil,
		},
		{
			name: "not existing call id",
			req: &api.GetCallRequest{
				CallId: "0000",
			},
			want:    nil,
			wantErr: service.ErrCallIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.GetCall(ctx, tt.req)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}

}

func TestGetLeaderCalls(t *testing.T) {
	ctx := context.Background()

	srv := service.New(zap.NewNop(), connP)

	tests := []struct {
		name    string
		req     *api.GetLeaderCallsRequest
		want    *api.GetLeaderCallsResponse
		wantErr error
	}{
		{
			name: "successful get leader calls",
			req: &api.GetLeaderCallsRequest{
				LeaderId: "leader2",
			},
			want: &api.GetLeaderCallsResponse{
				Calls: []*api.Call{
					{
						CallId:    "2222",
						UserId:    "user2",
						LeaderId:  "leader2",
						Title:     "title2",
						Status:    "status2",
						Feedback:  "feedback2",
						StartTime: "02:02:02",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "two calls",
			req: &api.GetLeaderCallsRequest{
				LeaderId: "leader3",
			},
			want: &api.GetLeaderCallsResponse{
				Calls: []*api.Call{
					{
						CallId:    "3333",
						UserId:    "user3",
						LeaderId:  "leader3",
						Title:     "title3",
						Status:    "status3",
						Feedback:  "feedback3",
						StartTime: "03:03:03",
					},
					{
						CallId:    "4444",
						UserId:    "user4",
						LeaderId:  "leader3",
						Title:     "title4",
						Status:    "status4",
						Feedback:  "feedback4",
						StartTime: "04:04:04",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "not existing call id",
			req: &api.GetLeaderCallsRequest{
				LeaderId: "leader0",
			},
			want:    nil,
			wantErr: service.ErrLeaderIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.GetLeaderCalls(ctx, tt.req)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.want)
		})
	}
}
