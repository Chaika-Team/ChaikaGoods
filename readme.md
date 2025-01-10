# ChaikaGoods

ChaikaGoods - это микросервис для управления товарами и шаблонами товаров, предназначенный для использования в проекте
Chaika. Сервис предоставляет функциональность для добавления, обновления, удаления и получения информации о товарах и
шаблонах.

## Основные возможности

- **Управление товарами**: добавление, обновление, удаление и получение информации о товарах.
- **Управление шаблонными списками товаров**: создание шаблонов, добавление товаров в шаблоны, удаление шаблонов и получение
  информации о содержимом шаблонов.
- **Поиск по товарам и шаблонам**: получение списка всех товаров или шаблонов, а также поиск по конкретным критериям.

## Технологии

ChaikaGoods использует следующий стек технологий:

- **Go** - язык программирования.
- **PostgreSQL** - система управления базами данных.
- **pgx** - библиотека для взаимодействия с PostgreSQL из Go.
- **Docker** - платформа для контейнеризации приложений.

## Настройка и запуск

Для работы с микросервисом необходимо выполнить следующие шаги:

### Предварительные требования

Убедитесь, что у вас установлены:

- Go (версия 1.15 или выше)
- Docker
- PostgreSQL (можно использовать Docker контейнер)

### Конфигурация

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/ChaikaGoods.git
   cd ChaikaGoods
   ```

2. Настройте файл конфигурации `config.yml`, указав параметры подключения к базе данных.

Пример файла конфигурации `config.yml`:

```yaml

log:
  level: info # debug, info, warn, error
listen:
  bind_ip: "127.0.0.1"
  port: 8080
storage:
  type: postgres
  host: 127.0.0.1
  port: 5432
  database: postgres
  user: postgres
  password: postgres
  max_conns: 10
  min_conns: 2
  health_check_period: 30s

```

### Запуск

Запустите микросервис локально:

```bash
go run main.go
```

### Запуск с использованием Docker Compose

Вы можете собрать Docker Compose и запустить микросервис с помощью следующих команд:

```bash
docker-compose build
docker-compose up
```

## Тестирование

Для запуска тестов выполните:

```bash
go test ./...
```

## Документация API

In progress...

## Лицензия

In progress...

## Контакты

Если у вас возникли вопросы или предложения, пожалуйста, не стесняйтесь обращаться ко мне
через [Issues](https://github.com/Chaika-Team/ChaikaGoods/issues) на GitHub.
