# Микросервисы — система заказов и уведомлений

Проект на Go с двумя сервисами: Order Service и Notification Service. Они взаимодействуют через Redis и PostgreSQL по событийной схеме.

## Что внутри

- **Order Service** — REST API для работы с заказами
- **Notification Service** — асинхронная обработка событий
- **PostgreSQL** — хранение данных
- **Redis** — доставка событий между сервисами

## Архитектура

```text
Order Service  --->  Redis  --->  Notification Service
      │                                │
      └──────── PostgreSQL ───────────┘
```

## Структура проекта

```
├── order-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── ...
├── notification-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── ...
├── docker-compose.yml
├── README.md
├── QUICKSTART.md
├── LICENSE
└── .gitignore
```

## Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.21+

### Запуск проекта

```bash
# Клонируйте репозиторий
git clone https://github.com/Alliazix/go-microservices.git
cd go-microservices

# Запустите все сервисы через Docker Compose
docker-compose up -d

# Проверьте логи
docker-compose logs -f

# Остановите сервисы
docker-compose down
```

### Локальное развертывание

```bash
# Запустите PostgreSQL и Redis
docker-compose up -d postgres redis

# Запустите Order Service
cd order-service
go run main.go

# В новом терминале запустите Notification Service
cd notification-service
go run main.go
```

## Тестирование API

```bash
# Создание заказа
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "items": [...], "total": 100}'

# Получение заказов
curl http://localhost:8080/orders
```

## Лицензия

MIT
