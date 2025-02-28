package postgresql_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient для тестирования без реального PostgreSQL
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	args := m.Called(append([]interface{}{ctx, sql}, arguments...)...)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *MockClient) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	mockArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return mockArgs.Get(0).(pgx.Rows), mockArgs.Error(1)
}

func (m *MockClient) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	mockArgs := m.Called(append([]interface{}{ctx, sql}, args...)...)
	return mockArgs.Get(0).(pgx.Row)
}

func (m *MockClient) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	mockArgs := m.Called(ctx, b)
	return mockArgs.Get(0).(pgx.BatchResults)
}

func (m *MockClient) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockClient) Close() {
	m.Called()
}

// Тест Exec()
func TestPGClient_Exec(t *testing.T) {
	mockClient := new(MockClient)
	ctx := context.Background()
	sql := "UPDATE table SET column = $1 WHERE id = $2"
	args := []any{"value", 1}
	tag := pgconn.NewCommandTag("UPDATE 1")

	mockClient.On("Exec", ctx, sql, args[0], args[1]).Return(tag, nil).Once()

	result, err := mockClient.Exec(ctx, sql, args...)

	assert.NoError(t, err)
	assert.Equal(t, tag, result)

	mockClient.AssertExpectations(t)
}

// // Тест Query() с pgxmock
// func TestPGClient_Query(t *testing.T) {
// 	mockClient := new(MockClient)
// 	mockDB, err := pgxmock.NewPool()
// 	assert.NoError(t, err)

// 	defer mockDB.Close()

// 	ctx := context.Background()
// 	sql := "SELECT * FROM table WHERE id = $1"
// 	args := []any{1}

// 	rows := pgxmock.NewRows([]string{"id", "name"}).AddRow(1, "Test")
// 	mockDB.ExpectQuery(sql).WithArgs(args...).WillReturnRows(rows)

// 	mockClient.On("Query", ctx, sql, args[0]).Return(rows, nil).Once()

// 	result, err := mockClient.Query(ctx, sql, args...)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)

// 	mockClient.AssertExpectations(t)
// }

// Тест QueryRow() с pgxmock
func TestPGClient_QueryRow(t *testing.T) {
	mockClient := new(MockClient)
	mockDB, err := pgxmock.NewPool()
	assert.NoError(t, err)

	defer mockDB.Close()

	ctx := context.Background()
	sql := "SELECT column FROM table WHERE id = $1"
	args := []any{1}

	row := mockDB.QueryRow(ctx, sql, args...)
	mockClient.On("QueryRow", ctx, sql, args[0]).Return(row).Once()

	result := mockClient.QueryRow(ctx, sql, args...)

	assert.NotNil(t, result)

	mockClient.AssertExpectations(t)
}

// Тест SendBatch()
// func TestPGClient_SendBatch(t *testing.T) {
// 	mockClient := new(MockClient)
// 	ctx := context.Background()
// 	batch := new(pgx.Batch)
// 	var mockBatchResults pgx.BatchResults // Заглушка

// 	mockClient.On("SendBatch", ctx, batch).Return(mockBatchResults).Once()

// 	results := mockClient.SendBatch(ctx, batch)

// 	assert.Equal(t, mockBatchResults, results)

// 	mockClient.AssertExpectations(t)
// }

// // Тест Begin()
// func TestPGClient_Begin(t *testing.T) {
// 	mockClient := new(MockClient)
// 	ctx := context.Background()
// 	var mockTx pgx.Tx // Заглушка

// 	mockClient.On("Begin", ctx).Return(mockTx, nil).Once()

// 	tx, err := mockClient.Begin(ctx)

// 	assert.NoError(t, err)
// 	assert.Equal(t, mockTx, tx)

// 	mockClient.AssertExpectations(t)
// }

// Тест Close()
func TestPGClient_Close(t *testing.T) {
	mockClient := new(MockClient)

	mockClient.On("Close").Return().Once()

	mockClient.Close()

	mockClient.AssertExpectations(t)
}
