#!/bin/bash

# Определяем путь к рабочей директории
WORK_DIR="$(dirname "$0")/.."

# Переходим в рабочую директорию
cd "$WORK_DIR" || { echo "Не удалось перейти в рабочую директорию $WORK_DIR"; exit 1; }

# Переменная для хранения пути к конфигурационному файлу
MOCKERY_CONFIG=".mockery.yaml"

# Проверяем наличие mockery в системе
if ! command -v mockery &> /dev/null; then
    echo "mockery не найден. Пожалуйста, установите его:"
    echo "  go install github.com/vektra/mockery/v2@latest"
    exit 1
fi

# Проверяем наличие конфигурационного файла
if [ ! -f "$MOCKERY_CONFIG" ]; then
    echo "Файл конфигурации $MOCKERY_CONFIG не найден. Пожалуйста, проверьте его наличие."
    exit 1
fi

# Запускаем mockery с указанным конфигурационным файлом
echo "Запускаем mockery с конфигурацией $MOCKERY_CONFIG..."
mockery --all --config="$MOCKERY_CONFIG"

# Проверяем статус выполнения
if [ $? -eq 0 ]; then
    echo "Генерация моков завершена успешно!"
else
    echo "Ошибка при генерации моков. Проверьте вывод выше."
    exit 1
fi
