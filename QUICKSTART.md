# Quick Start Guide

## 🚀 5-Minute Startup

### Prerequisites
- Docker & Docker Compose installed
- Git

### Step 1: Clone and navigate
```bash
cd c:\Users\Asus\Desktop\go_project
```

### Step 2: Start services
```bash
docker-compose up --build
```

Wait for all services to be healthy (you'll see "Notification Service started successfully" and "Order Service started successfully")

### Step 3: Test it out

#### Terminal 1: Create an order
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

You'll get back:
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

#### Terminal 2: Get notifications
```bash
curl http://localhost:8081/api/notifications | jq .
```

You should see a notification was created for the order!

#### Terminal 3: Cancel the order
```bash
curl -X PATCH http://localhost:8080/api/orders/{order-id}/cancel
```

#### Terminal 4: Check notifications again
```bash
curl http://localhost:8081/api/notifications | jq .
```

Now you'll see a cancellation notification too!

## 📊 Service URLs

| Service | URL | Purpose |
|---------|-----|---------|
| **Order Service** | http://localhost:8080 | Create & manage orders |
| **Notification Service** | http://localhost:8081 | View notifications |
| **PostgreSQL** | localhost:5432 | Database |
| **Redis** | localhost:6379 | Message Queue |

## 🛠️ Useful Commands

### Using Make (if available)
```bash
make up              # Start services
make down            # Stop services
make logs            # View all logs
make seed-order      # Create test order
make get-orders      # List all orders
make get-notifications  # List all notifications
```

### Using Docker Compose
```bash
docker-compose up -d           # Start in background
docker-compose logs -f         # View logs
docker-compose down            # Stop services
docker-compose ps              # List running services
```

### Health Checks
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
```

## 📝 Example Workflow

```bash
# 1. Start services
docker-compose up -d

# 2. Wait for services to be ready (check logs)
docker-compose logs

# 3. Create your first order
RESPONSE=$(curl -s -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "Jane Smith",
    "product": "Mouse",
    "quantity": 2,
    "price": 25.99
  }')

ORDER_ID=$(echo $RESPONSE | jq -r '.data.id')
echo "Created order: $ORDER_ID"

# 4. Check it was created
curl http://localhost:8080/api/orders/$ORDER_ID

# 5. Wait a moment for notification to process
sleep 1

# 6. Check notifications
curl http://localhost:8081/api/notifications

# 7. Cancel the order
curl -X PATCH http://localhost:8080/api/orders/$ORDER_ID/cancel

# 8. See the cancellation notification
curl http://localhost:8081/api/notifications
```

## 🐛 Troubleshooting

### Services won't start?
```bash
# Check logs
docker-compose logs

# Rebuild images
docker-compose build --no-cache

# Start fresh
docker-compose down -v
docker-compose up --build
```

### Port already in use?
```bash
# Find process using port 8080
# Windows:
netstat -ano | findstr :8080

# Kill the process
taskkill /PID <PID> /F
```

### Database connection issues?
```bash
# Check PostgreSQL
docker-compose logs postgres

# Ping database
docker-compose exec postgres pg_isready

# Connect to database
docker-compose exec postgres psql -U postgres -d microservices
```

### Redis connection issues?
```bash
# Check Redis
docker-compose logs redis

# Ping Redis
docker-compose exec redis redis-cli ping
```

## 📚 Next Steps

1. **Read the documentation**
   - See [README.md](README.md) for full documentation
   - Check [API_TESTING.md](API_TESTING.md) for all API endpoints
   - Review [DEVELOPMENT.md](DEVELOPMENT.md) for local development

2. **Explore the code**
   - Look at the service implementations
   - Review the database schema
   - Check the Redis event structure

3. **Modify and extend**
   - Add new endpoints
   - Create additional services
   - Implement new features

4. **Deploy to production**
   - Push to GitHub
   - Set up CI/CD
   - Deploy to cloud (AWS/GCP/Azure)

## 💡 Portfolio Tips

This project demonstrates:
- ✅ **Microservices Architecture** - Two independent services
- ✅ **Event-Driven Design** - Redis for async communication
- ✅ **Database Design** - PostgreSQL with proper schema
- ✅ **REST API** - Standard HTTP endpoints
- ✅ **Docker** - Container orchestration
- ✅ **Clean Code** - Well-organized structure
- ✅ **Error Handling** - Proper error management
- ✅ **Logging** - Structured logging

**Pro tip**: Add this to your GitHub with good documentation and you've got a solid portfolio project!

---

**Happy coding! 🎉**

Need help? Check the docs or open an issue on GitHub!
