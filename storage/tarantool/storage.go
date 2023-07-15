package tarantool

import (
	"fmt"
	"io/fs"

	"github.com/tarantool/go-tarantool/v2"

	"github.com/vstdy/otus-highload/storage"
	"github.com/vstdy/otus-highload/storage/tarantool/migrations"
)

const (
	serviceName = "tarantool"

	dbTableLoggingKey     = "db-table"
	dbOperationLoggingKey = "db-operation"
)

var _ storage.IMessageStorage = (*Storage)(nil)

type (
	// Storage keeps tarantool storage dependencies.
	Storage struct {
		config Config
		conn   *tarantool.Connection
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

// NewStorage creates a new tarantool Storage with custom options.
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

	conn, err := tarantool.Connect(st.config.Address, st.config.ToTarantoolOpts())
	if err != nil {
		return nil, fmt.Errorf("connection for URL (%s) failed: %w", st.config.Address, err)
	}

	st.conn = conn

	if _, err = st.conn.Do(tarantool.NewPingRequest()).Get(); err != nil {
		return nil, fmt.Errorf("ping for address (%s) failed: %w", st.config.Address, err)
	}

	return st, nil
}

// Close closes DB connection.
func (st *Storage) Close() error {
	return st.conn.Close()
}

// MigrateUp applies all available migrations.
func (st *Storage) MigrateUp() error {
	fsys := migrations.GetMigrations()
	luaMigrationFiles, err := fs.Glob(fsys, "*.lua")
	if err != nil {
		return fmt.Errorf("getting lua files: %w", err)
	}

	type Tuple struct {
		// Instruct msgpack to pack this struct as array, so no custom packer
		// is needed.
		_msgpack struct{} `msgpack:",asArray"` //nolint: structcheck,unused
		Id       uint
		Msg      string
		Name     string
	}

	var t []Tuple
	for _, file := range luaMigrationFiles {
		expr, errN := fsys.ReadFile(file)
		if errN != nil {
			return fmt.Errorf("reading file: %w", errN)
		}

		req := tarantool.NewEvalRequest(string(expr))
		//resp, errN := st.conn.Do(req).Get()
		errN = st.conn.Do(req).GetTyped(&t)
		if errN != nil {
			return fmt.Errorf("performing migrations: %w", errN)
		}
		//_ = resp

		fmt.Println("Error", err)
		fmt.Println("Data", t)
	}

	return nil
}
