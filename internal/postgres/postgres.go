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

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" envDefault:"5432"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" envDefault:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" envDefault:"1234"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" envDefault:"postgres"`
}

func New(ctx context.Context, config Config) (*pgx.Conn, error) {
	conString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	conn, err := pgx.Connect(ctx, conString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %w", err)
	}

	migration, err := migrate.New("file://./storage/migrations", conString)
	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}

	return conn, nil
}
