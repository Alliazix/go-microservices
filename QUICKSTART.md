# Руководство быстрого старта

## 🚀 Запуск за 5 минут

### Требования
- Docker и Docker Compose установлены
- Git

### Шаг 1: Клонируйте репозиторий
```bash
git clone https://github.com/Alliazix/go-microservices.git
cd go-microservices
```

### Шаг 2: Запустите сервисы
```bash
docker-compose up --build
```

Ожидайте, пока все сервисы будут готовы (вы увидите сообщения "Notification Service started successfully" и "Order Service started successfully")

### Шаг 3: Протестируйте

#### Терминал 1: Создайте заказ
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Doe",
    "product": "Laptop",
    "quantity": 1,
    "price": 999.99
  }'
```

Вы получите ответ:
```json
{
  "success": true,
  "data": {
    "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    "customer_name": "John Doe",
    "product": "Laptop",
    ...
  }
}
```

#### Терминал 2: Получите уведомления
```bash
curl http://localhost:8081/api/notifications | jq .
```

Вы должны увидеть уведомление о созданном заказе!

#### Терминал 3: Отмените заказ
```bash
curl -X PATCH http://localhost:8080/api/orders/{order-id}/cancel
```

#### Терминал 4: Проверьте уведомления снова
```bash
curl http://localhost:8081/api/notifications | jq .
```

Теперь вы увидите также уведомление об отмене!

## 📊 URLs сервисов

| Сервис | URL | Назначение |
|---------|-----|---------|
| **Order Service** | http://localhost:8080 | Создание и управление заказами |
| **Notification Service** | http://localhost:8081 | Просмотр уведомлений |
| **PostgreSQL** | localhost:5432 | База данных |
| **Redis** | localhost:6379 | Очередь сообщений |

## 🛠️ Полезные команды

### Используя Make (если установлен)
```bash
make up              # Запустить сервисы
make down            # Остановить сервисы
make logs            # Просмотр логов
make seed-order      # Создать тестовый заказ
make get-orders      # Список всех заказов
make get-notifications  # Список всех уведомлений
```

### Используя Docker Compose
```bash
docker-compose up -d           # Запустить в фоне
docker-compose logs -f         # Просмотр логов
docker-compose down            # Остановить сервисы
docker-compose ps              # Список запущенных сервисов
```

### Проверка здоровья
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
```

## 📝 Пример рабочего процесса

```bash
# 1. Запустите сервисы
docker-compose up -d

# 2. Дождитесь готовности (проверьте логи)
docker-compose logs

# 3. Создайте первый заказ
RESPONSE=$(curl -s -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Jane Smith",
    "product": "Mouse",
    "quantity": 2,
    "price": 25.99
  }')

ORDER_ID=$(echo $RESPONSE | jq -r '.data.id')
echo "Создан заказ: $ORDER_ID"

# 4. Проверьте создание
curl http://localhost:8080/api/orders/$ORDER_ID

# 5. Подождите обработки уведомления
sleep 1

# 6. Проверьте уведомления
curl http://localhost:8081/api/notifications

# 7. Отмените заказ
curl -X PATCH http://localhost:8080/api/orders/$ORDER_ID/cancel

# 8. Посмотрите уведомление об отмене
curl http://localhost:8081/api/notifications
```

## 🐛 Решение проблем

### Сервисы не запускаются?
```bash
# Проверьте логи
docker-compose logs

# Пересоберите образы
docker-compose build --no-cache

# Начните заново
docker-compose down -v
docker-compose up --build
```

### Порт уже используется?
```bash
# Найти процесс, использующий порт 8080
# Windows:
netstat -ano | findstr :8080

# Завершить процесс
taskkill /PID <PID> /F

# macOS/Linux:
lsof -i :8080
kill -9 <PID>
```

### Проблемы с подключением БД?
```bash
# Проверьте PostgreSQL
docker-compose logs postgres

# Проверьте доступность БД
docker-compose exec postgres pg_isready

# Подключитесь к БД
docker-compose exec postgres psql -U postgres -d microservices
```

### Проблемы с Redis?
```bash
# Проверьте Redis
docker-compose logs redis

# Проверьте Redis
docker-compose exec redis redis-cli ping
```

## 📚 Следующие шаги

1. **Прочитайте документацию**
   - Смотрите [README.md](README.md) для полной документации
   - Проверьте все API endpoints
   - Просмотрите локальную разработку

2. **Исследуйте код**
   - Смотрите реализацию сервисов
   - Изучите схему базы данных
   - Проверьте структуру Redis событий

3. **Модифицируйте и расширяйте**
   - Добавьте новые endpoints
   - Создайте дополнительные сервисы
   - Реализуйте новые функции

4. **Разверните в production**
   - Отправьте на GitHub
   - Настройте CI/CD
   - Разверните на облако (AWS/GCP/Azure)

## 💡 Советы для портфолио

Этот проект демонстрирует:
- ✅ **Архитектура микросервисов** — два независимых сервиса
- ✅ **Event-Driven дизайн** — Redis для асинхронной коммуникации
- ✅ **Дизайн БД** — PostgreSQL с правильной схемой
- ✅ **REST API** — стандартные HTTP endpoints
- ✅ **Docker** — оркестровка контейнеров
- ✅ **Чистый код** — хорошо организованная структура
- ✅ **Обработка ошибок** — правильное управление ошибками
- ✅ **Логирование** — структурированное логирование

**Совет**: Добавьте этот проект на GitHub с хорошей документацией и у вас будет отличный портфолио!

---

**Удачи в кодировании! 🎉**

Нужна помощь? Проверьте документацию или откройте issue на GitHub!
