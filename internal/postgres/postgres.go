package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" envDefault:"5432"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" envDefault:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" envDefault:"1234"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" envDefault:"postgres"`
	MaxConn  int32  `yaml:"POSTGRES_MAX_CONN" env:"POSTGRES_MAX_CONN" env-default:"10"`
	MinConn  int32  `yaml:"POSTGRES_MIN_CONN" env:"POSTGRES_MIN_CONN" env-default:"5"`
}

func New(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	cfgForPool := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.MaxConn,
		config.MinConn,
	)

	cfgForMigration := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	conn, err := pgxpool.New(ctx, cfgForPool)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %w", err)
	}

	migration, err := migrate.New("file://./storage/migrations", cfgForMigration)
	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return conn, nil
}
