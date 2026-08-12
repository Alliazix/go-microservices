# Руководство быстрого старта

## 🚀 Запуск за 5 минут

Я создал этот проект как демонстрацию своих навыков в разработке микросервисной архитектуры на Go. Вот как его запустить и посмотреть, как всё работает.

### Что вам понадобится

- Docker и Docker Compose (чтобы не устанавливать PostgreSQL и Redis вручную)
- Git (чтобы клонировать репозиторий)
- Curl или Postman (для тестирования API)

### Быстрый старт (3 шага)

#### Шаг 1: Клонируйте и перейдите в папку проекта

```bash
git clone https://github.com/Alliazix/go-microservices.git
cd go-microservices
```

#### Шаг 2: Запустите все сервисы

```bash
docker-compose up --build
```

Дождитесь, пока вы увидите в логах сообщения вроде:
```
order-service_1        | Order Service started on :8080
notification-service_1 | Notification Service started on :8081
```

#### Шаг 3: Протестируйте API

Откройте несколько терминалов и попробуйте следующие команды:

### Тестирование в действии

#### 1️⃣ Создайте заказ (Терминал 1)

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "product": "Ноутбук",
    "quantity": 1,
    "price": 999.99
  }'
```

Вы получите ответ с ID заказа:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "customer_name": "John Doe",
    "product": "Ноутбук",
    "quantity": 1,
    "price": 999.99,
    "status": "created",
    "created_at": "2026-08-12T15:30:00Z"
  }
}
```

#### 2️⃣ Получите список уведомлений (Терминал 2)

Сразу же откройте другой терминал и проверьте уведомления:

```bash
curl http://localhost:8081/api/notifications | jq .
```

Вы увидите, что Notification Service **асинхронно** обработал событие создания заказа через Redis! Вот в чём вся прелесть микросервисной архитектуры 😎

```json
{
  "success": true,
  "data": [
    {
      "id": "notification-001",
      "order_id": "550e8400-e29b-41d4-a716-446655440000",
      "message": "Заказ успешно создан",
      "type": "order_created",
      "created_at": "2026-08-12T15:30:01Z"
    }
  ]
}
```

#### 3️⃣ Отмените заказ (Терминал 1)

```bash
curl -X PATCH http://localhost:8080/api/orders/550e8400-e29b-41d4-a716-446655440000/cancel
```

#### 4️⃣ Снова проверьте уведомления (Терминал 2)

```bash
curl http://localhost:8081/api/notifications | jq .
```

Теперь вы увидите **два** уведомления — одно о создании, одно об отмене. Вот так event-driven архитектура работает на практике!

## 📊 Где что находится

| Сервис | URL | Что там |
|---------|-----|---------|
| **Order Service** | http://localhost:8080 | REST API для управления заказами |
| **Notification Service** | http://localhost:8081 | API для просмотра уведомлений |
| **PostgreSQL** | localhost:5432 | БД (пользователь: postgres, пароль: postgres) |
| **Redis** | localhost:6379 | Message Queue для общения между сервисами |

## 🛠️ Часто используемые команды

### Docker Compose команды

```bash
# Запустить сервисы в фоне
docker-compose up -d

# Посмотреть логи всех сервисов
docker-compose logs -f

# Посмотреть логи конкретного сервиса
docker-compose logs -f order-service
docker-compose logs -f notification-service

# Остановить всё
docker-compose down

# Посмотреть, какие контейнеры запущены
docker-compose ps

# Пересобрать образы и запустить заново
docker-compose up --build
```

### Проверка здоровья сервисов

```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
```

Оба должны вернуть 200 OK.

## 📝 Полный пример работы

Вот я запускаю проект, создаю заказ, проверяю, что уведомление пришло, отменяю заказ:

```bash
# Запустить всё
docker-compose up -d

# Дождаться загрузки (примерно 5 секунд)
sleep 5

# Создаём заказ
RESPONSE=$(curl -s -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Jane Smith",
    "product": "Мышка",
    "quantity": 2,
    "price": 25.99
  }')

# Берём ID заказа из ответа
ORDER_ID=$(echo $RESPONSE | jq -r '.data.id')
echo "✅ Заказ создан с ID: $ORDER_ID"

# Проверяем, создался ли заказ
echo "📦 Информация о заказе:"
curl http://localhost:8080/api/orders/$ORDER_ID | jq .

# Подождём, пока уведомление обработается
sleep 1

# Проверяем уведомления
echo "📢 Уведомления:"
curl http://localhost:8081/api/notifications | jq .

# Отменяем заказ
echo "❌ Отменяю заказ..."
curl -X PATCH http://localhost:8080/api/orders/$ORDER_ID/cancel

# Снова проверяем уведомления (теперь их два!)
echo "📢 Обновленные уведомления:"
curl http://localhost:8081/api/notifications | jq .

# Останавливаем
docker-compose down
```

## 🐛 Если что-то пошло не так

### Сервисы не запускаются

```bash
# Проверьте логи
docker-compose logs

# Если помогает, пересоберите всё с нуля
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```

### "Port already in use" — порт занят

Это значит, что на порте 8080 или 8081 уже что-то слушает.

```bash
# Windows:
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# macOS/Linux:
lsof -i :8080
kill -9 <PID>
```

### Проблемы с PostgreSQL

```bash
# Проверьте логи БД
docker-compose logs postgres

# Подключитесь к БД и посмотрите таблицы
docker-compose exec postgres psql -U postgres -d microservices -c "\dt"
```

### Проблемы с Redis

```bash
# Проверьте Redis
docker-compose logs redis

# Пингуем Redis
docker-compose exec redis redis-cli ping
```

## 💡 Что я здесь показываю

Этот проект демонстрирует, что я умею:

- ✅ **Микросервисная архитектура** — два независимых сервиса, каждый отвечает за свою часть
- ✅ **Event-Driven Design** — используется Redis Pub/Sub для асинхронной коммуникации между сервисами
- ✅ **REST API** — правильно спроектированные endpoints с правильными HTTP методами
- ✅ **PostgreSQL** — работаю с реальной БД, понимаю схемы и миграции
- ✅ **Docker & Docker Compose** — могу контейнеризировать приложения и управлять ими
- ✅ **Обработка ошибок** — приложение не падает при ошибках
- ✅ **Логирование** — в логах видно, что происходит в приложении
- ✅ **Чистая архитектура** — код хорошо организован, легко добавлять новые функции

## 🚀 Что дальше?

Если вам понравилось, вот что я планирую добавить:

1. **Аутентификация** — JWT токены для защиты API
2. **Тесты** — unit и integration тесты с >80% покрытием
3. **Метрики** — Prometheus и Grafana для мониторинга
4. **Кэширование** — Redis для кэширования часто запрашиваемых данных
5. **gRPC** — дополнительный протокол для межсервисной коммуникации
6. **Kubernetes** — развёртывание в k8s

## 📞 Вопросы?

Если у вас есть вопросы или вы нашли баг, [откройте issue](https://github.com/Alliazix/go-microservices/issues) на GitHub!

---

**Happy coding! 🎉**
