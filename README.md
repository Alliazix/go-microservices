# Microservices - Order & Notification System

Simple microservices project with two Go services that communicate via Redis and PostgreSQL. Good for portfolio and learning.

## What's Inside

- **Order Service** - REST API to create and manage orders (port 8080)
- **Notification Service** - Processes order events asynchronously via Redis (port 8081)
- **PostgreSQL** - Stores orders and notifications
- **Redis** - Event messaging between services

## Quick Start

### Start Everything
```bash
docker-compose up --build
```

Wait for services to start, then test:

### Create an Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John",
    "product": "Laptop",
    "quantity": 1,
    "price": 999.99
  }'
```

### Get All Orders
```bash
curl http://localhost:8080/api/orders
```

### Get Notifications
```bash
curl http://localhost:8081/api/notifications
```

### Cancel Order
```bash
curl -X PATCH http://localhost:8080/api/orders/{order-id}/cancel
```

## Project Structure

```
.
├── order-service/
│   ├── main.go         # Single file with all logic
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── notification-service/
│   ├── main.go         # Single file with all logic
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

## How It Works

1. Client creates an order via REST API
2. Order Service saves to PostgreSQL
3. Order Service publishes event to Redis
4. Notification Service listens to Redis
5. Notification Service creates notification record
6. Client can retrieve notifications

## Tech Stack

- Go 1.21
- PostgreSQL 15
- Redis 7
- Docker

## API Endpoints

### Order Service (8080)
- `POST /api/orders` - Create order
- `GET /api/orders` - List all orders
- `GET /api/orders/{id}` - Get order
- `PATCH /api/orders/{id}/cancel` - Cancel order
- `GET /health` - Health check

### Notification Service (8081)
- `GET /api/notifications` - List all notifications
- `GET /api/notifications/{id}` - Get notification
- `GET /health` - Health check

## Environment Variables

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=microservices
REDIS_URL=redis://localhost:6379
PORT=8080 (for order-service)
PORT=8081 (for notification-service)
```

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL or Docker
- Redis or Docker

### Local Setup

1. Start database and redis:
```bash
docker-compose up postgres redis
```

2. Set env variables:
```bash
export DB_HOST=localhost
export REDIS_URL=redis://localhost:6379
```

3. Run Order Service:
```bash
cd order-service
go run main.go
```

4. Run Notification Service (in another terminal):
```bash
cd notification-service
go run main.go
```

## Why This Project?

Good for learning:
- Microservices basics
- Event-driven architecture
- Go HTTP servers
- Database operations
- Docker & Docker Compose
- REST APIs
- Async messaging

Simple and realistic like junior developer code.

## License

MIT


## Project Overview

This project implements a distributed order processing system with two microservices:

- **Order Service**: REST API for creating and managing orders
- **Notification Service**: Asynchronous event processor for sending notifications

### Architecture

```
┌─────────────────┐
│  Order Service  │
│   (REST API)    │
│   Port: 8080    │
└────────┬────────┘
         │ Publishes events
         ▼
    ┌─────────┐
    │  Redis  │ ◄──────┐
    │ Message │        │ Subscribes to events
    │ Queue   │        │
    └─────────┘ ───────┘
         │
         │ Stores/Queries Data
         ▼
    ┌──────────────┐
    │  PostgreSQL  │
    └──────────────┘
         ▲
         │ Queries
         │
    ┌─────────────────────┐
    │ Notification Service│
    │   (Worker)          │
    │   Port: 8081        │
    └─────────────────────┘
```

## Tech Stack

- **Language**: Go 1.21+
- **Databases**: PostgreSQL, Redis
- **Messaging**: Redis Streams/Pub-Sub
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing + HTTP tests
- **Logging**: Structured logging

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Git

### Quick Start

1. **Clone and navigate:**
```bash
cd c:\Users\Asus\Desktop\go_project
```

2. **Start all services:**
```bash
docker-compose up --build
```

3. **Create an order (Order Service):**
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

4. **Check notifications (Notification Service):**
```bash
curl http://localhost:8081/api/notifications
```

## Project Structure

```
go_project/
├── docker-compose.yml           # Services orchestration
├── README.md                    # This file
├── .gitignore                   # Git ignore rules
├── order-service/               # Microservice 1
│   ├── main.go                  # Entry point
│   ├── go.mod                   # Go dependencies
│   ├── Dockerfile               # Container config
│   ├── cmd/
│   │   └── app/
│   │       └── main.go          # Application setup
│   └── internal/
│       ├── models/              # Data structures
│       ├── handlers/            # HTTP handlers
│       ├── service/             # Business logic
│       ├── db/                  # Database operations
│       ├── redis/               # Redis messaging
│       ├── config/              # Configuration
│       └── logger/              # Logging utilities
└── notification-service/        # Microservice 2
    ├── main.go                  # Entry point
    ├── go.mod                   # Go dependencies
    ├── Dockerfile               # Container config
    └── internal/
        ├── models/              # Data structures
        ├── service/             # Business logic
        ├── db/                  # Database operations
        ├── redis/               # Redis consumer
        ├── config/              # Configuration
        └── logger/              # Logging utilities
```

## Key Features

### Order Service
- ✅ RESTful API for order management
- ✅ Order validation
- ✅ Event publishing to Redis
- ✅ Database persistence
- ✅ Error handling & logging

### Notification Service
- ✅ Redis event subscription
- ✅ Asynchronous processing
- ✅ Notification persistence
- ✅ Status tracking
- ✅ Graceful shutdown

### Infrastructure
- ✅ Docker Compose for orchestration
- ✅ Health checks
- ✅ Environment-based configuration
- ✅ Structured logging
- ✅ Error handling

## API Endpoints

### Order Service (Port 8080)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/orders` | Create new order |
| GET | `/api/orders/{id}` | Get order details |
| GET | `/api/orders` | List all orders |
| PATCH | `/api/orders/{id}/cancel` | Cancel order |
| GET | `/health` | Health check |

### Notification Service (Port 8081)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/notifications` | List all notifications |
| GET | `/api/notifications/{id}` | Get notification details |
| GET | `/health` | Health check |

## Development

### Local Setup (without Docker)

```bash
# Install dependencies
go mod download

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=microservices
export REDIS_URL=redis://localhost:6379

# Run Order Service
cd order-service && go run main.go

# Run Notification Service (in another terminal)
cd notification-service && go run main.go
```

### Running Tests

```bash
# Order Service tests
cd order-service
go test ./... -v

# Notification Service tests
cd notification-service
go test ./... -v
```

## Database Schema

### Orders Table
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_name VARCHAR(255) NOT NULL,
    product VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Notifications Table
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP
);
```

## Communication Flow

1. **Order Creation**: Client sends order to Order Service
2. **Event Publishing**: Order Service publishes event to Redis
3. **Event Processing**: Notification Service subscribes to Redis events
4. **Notification Creation**: Notification Service creates notification record in DB
5. **Response**: Client can query both services for data

## Lessons Learned & Portfolio Value

This project demonstrates:

- 🏗️ **Microservices Architecture**: Separation of concerns
- 📡 **Async Communication**: Redis for decoupled services
- 💾 **Database Design**: Proper schema with constraints
- 🐳 **Containerization**: Docker best practices
- 📝 **Code Organization**: Clean architecture principles
- 🔄 **Error Handling**: Robust error management
- 🧪 **Testing**: Unit and integration tests
- 📊 **Logging**: Structured logging for debugging
- 🔌 **Configuration**: Environment-based config
- ⚡ **Performance**: Efficient async processing

## Scaling Considerations

- Redis Cluster for high availability
- PostgreSQL read replicas
- Load balancing with nginx
- Service mesh (Istio/Linkerd)
- Horizontal pod autoscaling (Kubernetes)

## Troubleshooting

### PostgreSQL connection issues
```bash
# Check PostgreSQL health
docker-compose exec postgres pg_isready
```

### Redis connection issues
```bash
# Check Redis health
docker-compose exec redis redis-cli ping
```

### View logs
```bash
docker-compose logs -f order-service
docker-compose logs -f notification-service
```

## Production Considerations

- Add API authentication (JWT)
- Implement rate limiting
- Add request validation middleware
- Set up distributed tracing
- Implement circuit breakers
- Add caching strategies
- Use message queue (RabbitMQ/Kafka) for scale

## License

MIT License - feel free to use for portfolio or learning

## Future Enhancements

- [ ] User authentication & authorization
- [ ] Webhook notifications
- [ ] Email/SMS integration
- [ ] Order analytics dashboard
- [ ] Payment processing integration
- [ ] Rate limiting & throttling
- [ ] API documentation (Swagger)
- [ ] Prometheus metrics
- [ ] Kubernetes deployment files

---

**Created for portfolio demonstration | Medium/Junior+ Level**
