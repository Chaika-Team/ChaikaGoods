package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// Тест 1 --------------------------------------------------------------------------------

// MockRow для мокирования ответа QueryRow
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

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

// Тест 2 --------------------------------------------------------------------------------

// MockRows для мокирования pgx.Rows
type MockRows struct {
	mock.Mock
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (m *MockRows) Conn() *pgx.Conn {
	return nil
}

func (m *MockRows) Close() {}
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}
func (m *MockRows) RawValues() [][]byte {
	return nil
}

// Исправление: добавляем метод Values()
func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

// Тест 4 --------------------------------------------------------------------------------

// MockResult для мокирования CommandTag (результата Exec)
// MockTx — мок для транзакций pgx.Tx
// MockTx — мок для транзакций pgx.Tx
type MockTx struct {
	mock.Mock
}

// QueryRow мокирует метод QueryRow для выполнения запроса с одним результатом.
func (m *MockTx) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	argsList := append([]interface{}{ctx, query}, args...)
	return m.Called(argsList...).Get(0).(pgx.Row)
}

// Exec мокирует выполнение SQL-запроса без возврата строк.
func (m *MockTx) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsList := append([]interface{}{ctx, query}, args...)
	res := m.Called(argsList...)
	return res.Get(0).(pgconn.CommandTag), res.Error(1)
}

// Commit мокирует фиксацию транзакции.
func (m *MockTx) Commit(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// Rollback мокирует откат транзакции.
func (m *MockTx) Rollback(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// Conn возвращает nil, так как мок не использует соединение напрямую.
func (m *MockTx) Conn() *pgx.Conn {
	return nil
}

// SendBatch мокирует отправку батч-запросов.
func (m *MockTx) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	args := m.Called(ctx, batch)
	return args.Get(0).(pgx.BatchResults)
}

// Prepare мокирует подготовку SQL-запроса в транзакции.
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	args := m.Called(ctx, name, sql)
	// Теперь возвращаем указатель
	stmtDesc, _ := args.Get(0).(pgconn.StatementDescription)
	return &stmtDesc, args.Error(1)
}

// CopyFrom мокирует вставку данных через COPY.
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columns []string, rowSrc pgx.CopyFromSource) (int64, error) {
	args := m.Called(ctx, tableName, columns, rowSrc)
	return args.Get(0).(int64), args.Error(1)
}

// BeginTx реализует метод BeginTx для транзакции (сам транзакции не создаёт, но удовлетворяет интерфейсу).
func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return m, nil
}

// LargeObjects возвращает nil, так как мок не использует LargeObjects.
func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

// Query мокирует выполнение SQL-запроса, возвращая строки.
func (m *MockTx) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	argsList := append([]interface{}{ctx, query}, args...)
	res := m.Called(argsList...)
	return res.Get(0).(pgx.Rows), res.Error(1)
}
