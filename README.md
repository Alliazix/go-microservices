# Микросервисы — система заказов и уведомлений

Проект на Go с двумя сервисами: Order Service и Notification Service. Они взаимодействуют через Redis и PostgreSQL по событийной схеме. Подходит для обучения, портфолио и демонстрации микросервисной архитектуры.

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
      Структура проэкта:
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
