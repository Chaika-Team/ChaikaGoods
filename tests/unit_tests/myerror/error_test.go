package myerror

import (
	"errors"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/stretchr/testify/assert"
)

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции Error
//   - Классы эквивалентности: в зависимости от наличия Cause, формат сообщения меняется
func TestErrorAppError(t *testing.T) {
	appErr := &myerr.AppError{
		Message: "Simple error",
		Cause:   nil,
	}
	expected := "Simple error"
	assert.Equal(t, expected, appErr.Error())

	underlying := errors.New("check env variables")
	appErr = &myerr.AppError{
		Message: "Validation error",
		Cause:   underlying,
	}
	expected = "Validation error: check env variables"
	assert.Equal(t, expected, appErr.Error())
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции Unwrap
//   - Классы эквивалентности: функция возвращает значение, которое ранее было передано в структуру
func TestUnwrapAppError(t *testing.T) {
	cause := errors.New("network was broken")
	appErr := &myerr.AppError{
		Message: "a problem was occured...",
		Cause:   cause,
	}
	assert.Equal(t, cause, appErr.Unwrap())

	appErr = &myerr.AppError{
		Message: "please check support!",
		Cause:   nil,
	}
	assert.Equal(t, nil, appErr.Unwrap())
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции New
//   - Проверяет, что все поля ошибки установлены корректно
func TestNewAppError(t *testing.T) {
	cause := errors.New("host is down")
	ctx := map[string]interface{}{"key": "value"}
	appErr := myerr.New(myerr.ErrorTypeNotFound, "Not found", cause, ctx)

	assert.Equal(t, myerr.ErrorTypeNotFound, appErr.Type)
	assert.Equal(t, "Not found", appErr.Message)
	assert.Equal(t, cause, appErr.Cause)
	assert.Equal(t, ctx, appErr.Context)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции IsAppError
//   - Если ошибка принадлежит классу AppError, то результат должен быть положительным
func TestIsAppErrorTrue(t *testing.T) {
	appErr := myerr.New(myerr.ErrorTypeNotFound, "Server not found", nil, nil)
	res, ok := myerr.IsAppError(appErr)
	assert.True(t, ok)
	assert.Equal(t, appErr, res)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции IsAppError
//   - Проверка, что ошибка не принадлежит классу AppError, поэтому возвращаемое значение false
func TestIsAppErrorFalse(t *testing.T) {
	stdErr := errors.New("Default error")
	res, ok := myerr.IsAppError(stdErr)
	assert.False(t, ok)
	assert.Nil(t, res)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции Wrap
//   - Проверяет, что происходит оборачивание ошибки при объявленных аргументах
func TestWrapNonNil(t *testing.T) {
	cause := errors.New("Default error")
	wrapped := myerr.Wrap(cause, myerr.ErrorTypeInternal, "this is wrapped error", nil)
	assert.NotNil(t, wrapped)
	assert.Equal(t, myerr.ErrorTypeInternal, wrapped.Type)
	assert.Equal(t, "this is wrapped error", wrapped.Message)
	assert.Equal(t, cause, wrapped.Cause)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции Wrap
//   - Проверяет случай, когда входная ошибка равна nil.
func TestWrapNil(t *testing.T) {
	wrapped := myerr.Wrap(nil, myerr.ErrorTypeInternal, "internal error!", nil)
	assert.Nil(t, wrapped)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для предопределённых конструкторов (NotFound, Validation, Internal, Duplicate, Unauthorized, Forbidden, Conflict, Unknown)
//   - Проверяем то, что для каждого конструктора всегда эквивалентна своя ошибка
func TestPredefinedErrorConstructors(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(string, error) *myerr.AppError
		errType     myerr.ErrorType
		message     string
	}{
		{"NotFound", myerr.NotFound, myerr.ErrorTypeNotFound, "not found"},
		{"Validation", myerr.Validation, myerr.ErrorTypeValidation, "validation error"},
		{"Internal", myerr.Internal, myerr.ErrorTypeInternal, "internal error"},
		{"Duplicate", myerr.Duplicate, myerr.ErrorTypeDuplicate, "duplicate"},
		{"Unauthorized", myerr.Unauthorized, myerr.ErrorTypeUnauthorized, "unauthorized"},
		{"Forbidden", myerr.Forbidden, myerr.ErrorTypeForbidden, "forbidden"},
		{"Conflict", myerr.Conflict, myerr.ErrorTypeConflict, "conflict"},
		{"Unknown", myerr.Unknown, myerr.ErrorTypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor(tt.message, nil)
			assert.Equal(t, tt.errType, err.Type)
			assert.Equal(t, tt.message, err.Message)
		})
	}
}

// Техника тест-дизайна: Причинно-следственный анализ
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции WithContext
//   - Анализируется случай, когда ошибка параметр уже является AppError.
func TestWithContextAppError(t *testing.T) {
	appErr := myerr.New(myerr.ErrorTypeValidation, "validation failed!", nil, nil)
	addCtx := map[string]interface{}{"field": "value"}
	updated := myerr.WithContext(appErr, addCtx)

	assert.NotNil(t, updated.Context)
	assert.Equal(t, "value", updated.Context["field"])
}

// Техника тест-дизайна: Причинно-следственный анализ
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции WithContext
//   - Анализируется случай, когда входная ошибка не является AppError.
func TestWithContextNonAppError(t *testing.T) {
	stdErr := errors.New("stdout error")
	addCtx := map[string]interface{}{"a": 1, "debug": true}
	updated := myerr.WithContext(stdErr, addCtx)
	assert.Equal(t, myerr.ErrorTypeUnknown, updated.Type)
	assert.Equal(t, "unknown error with context", updated.Message)
	assert.Equal(t, addCtx, updated.Context)
	assert.Equal(t, stdErr, updated.Cause)
}

// Техника тест-дизайна: Попарное тестирование
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции IsType
//   - Попарное тестирование: рассматриваются случаи, когда передаются комбинации ошибкок разного вида и разного типа
func TestIsType(t *testing.T) {
	appErr := myerr.New(myerr.ErrorTypeDuplicate, "Duplication error!", nil, nil)
	nonAppErr := errors.New("Duplication error")
	assert.True(t, myerr.IsType(appErr, myerr.ErrorTypeDuplicate))
	assert.False(t, myerr.IsType(appErr, myerr.ErrorTypeNotFound))
	assert.False(t, myerr.IsType(nonAppErr, myerr.ErrorTypeDuplicate))
	assert.False(t, myerr.IsType(nonAppErr, myerr.ErrorTypeNotFound))
}

// Техника тест-дизайна: Таблица принятия решений
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функций IsNotFound, IsValidation, IsInternal, IsDuplicate, IsUnauthorized, IsForbidden, IsConflict, IsUnknown.
func TestIsSpecificErrorTypes(t *testing.T) {
	tests := []struct {
		name                string
		constructor         func(string, error) *myerr.AppError
		checkFunc           func(error) bool
		negativeConstructor func(string, error) *myerr.AppError
	}{
		{
			name:                "IsNotFound",
			constructor:         myerr.NotFound,
			checkFunc:           myerr.IsNotFound,
			negativeConstructor: myerr.Internal,
		},
		{
			name:                "IsValidation",
			constructor:         myerr.Validation,
			checkFunc:           myerr.IsValidation,
			negativeConstructor: myerr.NotFound,
		},
		{
			name:                "IsInternal",
			constructor:         myerr.Internal,
			checkFunc:           myerr.IsInternal,
			negativeConstructor: myerr.NotFound,
		},
		{
			name:                "IsDuplicate",
			constructor:         myerr.Duplicate,
			checkFunc:           myerr.IsDuplicate,
			negativeConstructor: myerr.NotFound,
		},
		{
			name:                "IsUnauthorized",
			constructor:         myerr.Unauthorized,
			checkFunc:           myerr.IsUnauthorized,
			negativeConstructor: myerr.NotFound,
		},
		{
			name:                "IsForbidden",
			constructor:         myerr.Forbidden,
			checkFunc:           myerr.IsForbidden,
			negativeConstructor: myerr.NotFound,
		},
		{
			name:                "IsConflict",
			constructor:         myerr.Conflict,
			checkFunc:           myerr.IsConflict,
			negativeConstructor: myerr.NotFound,
		},
		{
			name:                "IsUnknown",
			constructor:         myerr.Unknown,
			checkFunc:           myerr.IsUnknown,
			negativeConstructor: myerr.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor("test", nil)
			assert.True(t, tt.checkFunc(err))
			otherErr := tt.negativeConstructor("other", nil)
			assert.False(t, tt.checkFunc(otherErr))
		})
	}
}
