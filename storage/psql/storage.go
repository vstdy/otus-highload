package psql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/vstdy/otus-highload/storage"
	"github.com/vstdy/otus-highload/storage/psql/migrations"
)

const (
	serviceName = "psql"

	dbTableLoggingKey     = "db-table"
	dbOperationLoggingKey = "db-operation"
)

var _ storage.IStorage = (*Storage)(nil)

type (
	// Storage keeps psql storage dependencies.
	Storage struct {
		config           Config
		masterConn       *pgxpool.Pool
		asyncReplicaConn *pgxpool.Pool
	}

	// StorageOption defines functional argument for Storage constructor.
	StorageOption func(st *Storage) error
)

// WithConfig overrides default Storage config.
func WithConfig(config Config) StorageOption {
	return func(st *Storage) error {
		st.config = config

		return nil
	}
}

// NewStorage creates a new psql Storage with custom options.
func NewStorage(opts ...StorageOption) (*Storage, error) {
	st := &Storage{
		config: NewDefaultConfig(),
	}
	for optIdx, opt := range opts {
		if err := opt(st); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := st.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	ctx := context.Background()

	masterConn, err := pgxpool.New(ctx, st.config.DSN)
	if err != nil {
		return nil, fmt.Errorf("connection for DSN (%s) failed: %w", st.config.DSN, err)
	}

	var asyncReplicaConn *pgxpool.Pool
	if st.config.AsyncReplicaDSN != "" {
		asyncReplicaConn, err = pgxpool.New(ctx, st.config.AsyncReplicaDSN)
		if err != nil {
			return nil, fmt.Errorf("connection for DSN (%s) failed: %w", st.config.DSN, err)
		}
	}

	st.masterConn = masterConn
	st.asyncReplicaConn = asyncReplicaConn

	if err = st.masterConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping for DSN (%s) failed: %w", st.config.DSN, err)
	}

	return st, nil
}

// Close closes DB connection.
func (st *Storage) Close() error {
	if st.masterConn == nil {
		return nil
	}

	st.masterConn.Close()

	return nil
}

// MigrateUp applies all available migrations.
func (st *Storage) MigrateUp() error {
	db, err := sql.Open("pgx", st.config.DSN)
	if err != nil {
		return fmt.Errorf("connection for DSN (%s) failed: %w", st.config.DSN, err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.GetMigrations())
	_ = goose.SetDialect("pgx")

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("performing migrations: %w", err)
	}

	return nil
}

// MigrateDown rolls back a single migration from the current version.
func (st *Storage) MigrateDown() error {
	db, err := sql.Open("pgx", st.config.DSN)
	if err != nil {
		return fmt.Errorf("connection for DSN (%s) failed: %w", st.config.DSN, err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.GetMigrations())
	_ = goose.SetDialect("pgx")

	err = goose.Down(db, ".")
	if err != nil {
		return fmt.Errorf("performing migrations: %w", err)
	}

	return nil
}
