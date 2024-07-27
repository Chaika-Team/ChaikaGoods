package utils

import (
	"fmt"
	"reflect"
	"sync"
)

var fieldCache sync.Map // Кэш для хранения информации о полях структур

// VerifyMapFields проверяет, что все ключи в map соответствуют полям заданной структуры T. Необязательно все поля должны быть представлены в map.
func VerifyMapFields[T any](data map[string]interface{}) error {
	typ := reflect.TypeOf((*T)(nil)).Elem()

	// Пытаемся получить список допустимых ключей из кэша
	if cached, ok := fieldCache.Load(typ); ok {
		validKeys := cached.(map[string]bool)
		return validateKeys(data, validKeys)
	}

	// Создаем и кэшируем список допустимых ключей, если они еще не были кэшированы
	validKeys := make(map[string]bool)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if commaIdx := len(jsonTag); commaIdx > 0 {
			jsonTag = jsonTag[:commaIdx]
		}
		if jsonTag != "" {
			validKeys[jsonTag] = true
		}
	}

	fieldCache.Store(typ, validKeys)
	return validateKeys(data, validKeys)
}

// validateKeys проверяет, что все ключи в map присутствуют в списке допустимых ключей.
func validateKeys(data map[string]interface{}, validKeys map[string]bool) error {
	for key := range data {
		if !validKeys[key] {
			return fmt.Errorf("invalid field: %s", key)
		}
	}
	return nil
}
