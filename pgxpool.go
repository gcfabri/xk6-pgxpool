package sql

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.k6.io/k6/js/modules"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/pgxpool", New())
}

type (
	// RootModule is the global module instance that will create module
	// instances for each Pool.
	RootModule struct{}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		// vu provides methods for accessing internal k6 objects for a Pool
		vu modules.VU
		// pool is the exported type
		pool *Pool
	}

	// KeyValue is a simple key-value pair.
	KeyValue map[string]interface{}
)

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface returning a new instance for each Pool.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:   vu,
		pool: &Pool{vu: vu},
	}
}

// Pool is the type for our custom API.
type Pool struct {
	vu modules.VU // provides methods for accessing internal k6 objects
}

// Open establishes a connection to the specified database type using the provided connection string.
func (p *Pool) Open(connString string, minConns int, maxConns int) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.MinConns = int32(minConns)
	config.MaxConns = int32(maxConns)

	pool, err := pgxpool.NewWithConfig(p.vu.Context(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func (p *Pool) Exec(pool *pgxpool.Pool, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return pool.Exec(p.vu.Context(), sql, args...)
}

func (p *Pool) Query(pool *pgxpool.Pool, sql string, args ...interface{}) ([]KeyValue, error) {
	rows, err := pool.Query(p.vu.Context(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []KeyValue
	columns := rows.FieldDescriptions()
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		result := make(KeyValue)
		for i, column := range columns {
			result[column.Name] = values[i]
		}
		results = append(results, result)
	}

	return results, nil
}

// Exports implements the modules.Instance interface and returns the exported types for the JS module.
func (m *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Default: m.pool,
	}
}
