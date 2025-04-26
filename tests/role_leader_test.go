package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5/pgxpool"

	"role-leader/internal/postgres"
)

func upDB() *pgxpool.Pool {
	ctx := context.Background()

	cfg := postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "root",
		Password: "1234",
		Database: "postgres",
		MaxConn:  10,
		MinConn:  5,
	}
	cfgForMigration := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	conn, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	defer conn.Close()
	fmt.Println(os.Getwd())
	//migrate.NewMigration()
	migration, err := migrate.New("file://./migrations-for-tests", cfgForMigration)
	if err != nil {
		log.Fatalf("unable to create migrations: %v", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("unable to run migrations: %v", err)
	}

	return conn
}

func TestCreateFeedback(t *testing.T) {
	//t.Parallel()
	//ctx := context.Background()

	conn := upDB()
	fmt.Println(conn)
	//
	//srv := service.New(nil, nil, conn)
	//tests := []struct {
	//	name string
	//	req  *api.CreateFeedbackRequest
	//	want pgx.Rows
	//}{
	//	{
	//		name: "test1",
	//		req: &api.CreateFeedbackRequest{
	//			CallId:  "1111",
	//			Message: "testMessage",
	//		},
	//		want: nil,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		srv.CreateFeedback(ctx, tt.req)
	//		rows, err := conn.Query(ctx, tt.req.CallId)
	//		if err != nil {
	//			t.Fatalf("unable to execute query: %v", err)
	//		}
	//		if reflect.DeepEqual(rows, tt.want) {
	//			t.Errorf("CreateFeedback() = %v, want %v", rows, tt.want)
	//		}
	//	})
	//}
}
