package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockClient — мок для тестирования GoodsPGRepository
type MockClient struct {
	mock.Mock
}

func (m *MockClient) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	argsList := append([]interface{}{ctx, query}, args...)
	return m.Called(argsList...).Get(0).(pgx.Row)
}

func (m *MockClient) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	argsList := append([]interface{}{ctx, query}, args...)
	res := m.Called(argsList...)
	return res.Get(0).(pgx.Rows), res.Error(1)
}

func (m *MockClient) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsList := append([]interface{}{ctx, query}, args...)
	res := m.Called(argsList...)
	return res.Get(0).(pgconn.CommandTag), res.Error(1)
}

// Добавляем метод Begin для транзакций
func (m *MockClient) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

// Добавляем метод SendBatch для поддержки батч-запросов
func (m *MockClient) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, b)
	return args.Get(0).(pgx.BatchResults)
}

// Добавляем Close, если используется в коде
func (m *MockClient) Close() {
	m.Called()
}

// Добавляем Ping, если используется
func (m *MockClient) Ping(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
