#!/bin/bash
set -euo pipefail

# Загрузить .env
export $(grep -v '^#' .env | xargs)

# Проверка, что база готова
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

# Применить дамп (на случай повторного запуска)
./scripts/db_restore.sh

# Запуск тестов
go test -v -tags=integration ../tests/integration_tests

