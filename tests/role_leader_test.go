package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"role-leader/internal/api"
	postgres2 "role-leader/internal/postgres"
	"role-leader/internal/service/grpcSrv"
)

func upDB2() {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	l, _ := zap.NewDevelopment()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "1234",
			"POSTGRES_USER":     "root",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})
	if err != nil {
		l.Error("Failed to start container", zap.Error(err))
	}

	host, err := container.Host(ctx)
	if err != nil {
		l.Error("Failed to get container", zap.Error(err))
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		l.Error("Failed to get mapped port", zap.Error(err))
	}

	//dsn := fmt.Sprintf("host=%s port=%s user=root password=1234 dbname=postgres sslmode=disable", host, port.Port())
	//dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool", "root", "1234", host, port, "postgres")
	cfg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool",
		"root",
		"1234",
		host,
		port.Port(),
		"postgres",
	)

	conn, err := pgxpool.New(ctx, cfg+"_max_conns=10&pool_min_conns=5")
	if err != nil {
		l.Error("Failed to create pool", zap.Error(err))
		return
	}

	if err := conn.Ping(ctx); err != nil {
		l.Error("Failed to ping pool", zap.Error(err))
		return
	} else {
		l.Info("Successfully connected to postgres")
	}

	cfgForMigration := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool",
		"root",
		"1234",
		host,
		port.Port(),
		"postgres",
	)

	migration, err := migrate.New("file://../storage/migrations-for-tests", cfgForMigration)
	if err != nil {
		l.Error("Failed to create migrations", zap.Error(err))
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error("Failed to run migrations", zap.Error(err))
	}

	//migration, err := migrate.New("file://../storage/migrations-for-tests", cfg)
	//if err != nil {
	//	t.Errorf("Failed to create migrations: %v", err)
	//}
	//
	//err = migration.Up()
	//if err != nil {
	//	t.Errorf("Failed to run migrations: %v", err)
	//}

	time.Sleep(50 * time.Second)

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			l.Error("Failed to terminate postgres container: %v", zap.Error(err))
		}
		conn.Close()
	}()
}

func TestCreateFeedback(t *testing.T) {

	ctx := context.Background()

	conn, err := upDB()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	srv := grpcSrv.New(zap.NewNop(), conn)

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
			wantErr: grpcSrv.ErrEmptyMessage,
		},
		{
			name: "not existing call id",
			req: &api.CreateFeedbackRequest{
				CallId:  "0000",
				Message: "aboba",
			},
			want:    &api.Call{},
			wantErr: grpcSrv.ErrCallIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "select * from schema_call.phone_call where call_id = $1"
			var got api.Call

			_, err := srv.CreateFeedback(ctx, tt.req)
			conn.QueryRow(ctx, q, tt.req.CallId).Scan(
				&got.CallId,
				&got.UserId,
				&got.LeaderId,
				&got.Title,
				&got.Status,
				&got.Feedback,
				&got.StartTime,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateFeedback() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("CreateFeedback() got = %v, want %v", &got, tt.want)
			}
		})
	}
}

func TestGetCall(t *testing.T) {
	ctx := context.Background()
	conn, err := upDB()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	srv := grpcSrv.New(zap.NewNop(), conn)

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
			wantErr: grpcSrv.ErrCallIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.GetCall(ctx, tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetCall() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCall() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestGetLeaderCalls(t *testing.T) {
	ctx := context.Background()
	conn, err := upDB()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	srv := grpcSrv.New(zap.NewNop(), conn)

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
			wantErr: grpcSrv.ErrLeaderIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.GetLeaderCalls(ctx, tt.req)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetLeaderCalls(), test name = %s: error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLeaderCalls() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func upDB() (*pgxpool.Pool, error) {
	ctx := context.Background()

	cfg := postgres2.Config{
		Host:     "localhost",
		Port:     "2345",
		Username: "root",
		Password: "1234",
		Database: "postgres",
		MaxConn:  10,
		MinConn:  5,
	}

	cfgForPool := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.MaxConn,
		cfg.MinConn,
	)

	cfgForMigration := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	conn, err := pgxpool.New(ctx, cfgForPool)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %w", err)
	}

	migration, err := migrate.New("file://../storage/migrations-for-tests", cfgForMigration)
	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return conn, nil
}
