// mocks/db_mocks.go
package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockPool represents a mock database pool that implements the necessary methods
type MockPool struct {
	mock.Mock
}

// Begin starts a transaction
func (m *MockPool) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

// Exec executes a query that doesn't return rows
func (m *MockPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return pgconn.CommandTag{}, args.Error(1)
	}
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

// Query executes a query that returns rows
func (m *MockPool) Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error) {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Rows), args.Error(1)
}

// QueryRow executes a query that is expected to return at most one row
func (m *MockPool) QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return &MockRow{}
	}
	return args.Get(0).(pgx.Row)
}

// Close closes the pool
func (m *MockPool) Close() {
	m.Called()
}

// MockTx represents a mock transaction
type MockTx struct {
	mock.Mock
}

// Begin starts a pseudo nested transaction
func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}

// Commit commits the transaction
func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Rollback rolls back the transaction
func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Exec executes a query that doesn't return rows
func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return pgconn.CommandTag{}, args.Error(1)
	}
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

// Query executes a query that returns rows
func (m *MockTx) Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error) {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Rows), args.Error(1)
}

// QueryRow executes a query that is expected to return at most one row
func (m *MockTx) QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row {
	args := m.Called(ctx, sql, arguments)
	if args.Get(0) == nil {
		return &MockRow{}
	}
	return args.Get(0).(pgx.Row)
}

// CopyFrom uses the PostgreSQL copy protocol to perform bulk data insertion
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columnNames, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

// SendBatch sends a batch of queries
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(pgx.BatchResults)
}

// LargeObjects returns a LargeObjects instance
func (m *MockTx) LargeObjects() pgx.LargeObjects {
	args := m.Called()
	return args.Get(0).(pgx.LargeObjects)
}

// Prepare creates a prepared statement
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pgconn.StatementDescription), args.Error(1)
}

// Conn returns the underlying connection
func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgx.Conn)
}

// MockRow represents a mock row
type MockRow struct {
	mock.Mock
	scanFunc func(dest ...interface{}) error
}

// Scan reads the values from the row
func (m *MockRow) Scan(dest ...interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	args := m.Called(dest...)
	return args.Error(0)
}

// MockRows represents mock rows
type MockRows struct {
	mock.Mock
}

// Close closes the rows
func (m *MockRows) Close() {
	m.Called()
}

// Err returns any error that occurred while reading
func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

// CommandTag returns the command tag from the query
func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

// FieldDescriptions returns the field descriptions
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]pgconn.FieldDescription)
}

// Next prepares the next row for reading
func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Scan reads the values from the current row
func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

// Values returns the decoded row values
func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

// RawValues returns the raw row values
func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([][]byte)
}

// Conn returns the underlying connection
func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgx.Conn)
}
