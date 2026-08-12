package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Order represents an order
type Order struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"customer_name"`
	Product      string    `json:"product"`
	Quantity     int       `json:"quantity"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	CustomerName string  `json:"customer_name"`
	Product      string  `json:"product"`
	Quantity     int     `json:"quantity"`
	Price        float64 `json:"price"`
}

type OrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Order `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Global variables
var (
	db          *sql.DB
	redisClient *redis.Client
)

func main() {
	// Initialize database
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "microservices"),
	)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✓ Connected to database")

	// Initialize database schema
	if err = initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	redisClient = redis.NewClient(opt)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("✓ Connected to Redis")

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/orders", createOrderHandler).Methods("POST")
	router.HandleFunc("/api/orders", listOrdersHandler).Methods("GET")
	router.HandleFunc("/api/orders/{id}", getOrderHandler).Methods("GET")
	router.HandleFunc("/api/orders/{id}/cancel", cancelOrderHandler).Methods("PATCH")

	port := getEnv("PORT", "8080")
	log.Printf("✓ Order Service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		customer_name VARCHAR(255) NOT NULL,
		product VARCHAR(255) NOT NULL,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		price DECIMAL(10, 2) NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
	CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
	`

	_, err := db.Exec(query)
	return err
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if req.CustomerName == "" || req.Product == "" || req.Quantity <= 0 || req.Price <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid order data")
		return
	}

	// Create order
	var order Order
	query := `
		INSERT INTO orders (customer_name, product, quantity, price, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, customer_name, product, quantity, price, status, created_at
	`

	err := db.QueryRow(query, req.CustomerName, req.Product, req.Quantity, req.Price, "pending").
		Scan(&order.ID, &order.CustomerName, &order.Product, &order.Quantity, &order.Price, &order.Status, &order.CreatedAt)

	if err != nil {
		log.Println("Error creating order:", err)
		respondError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// Publish event to Redis
	event := map[string]interface{}{
		"id":            order.ID,
		"customer_name": order.CustomerName,
		"product":       order.Product,
		"quantity":      order.Quantity,
		"price":         order.Price,
		"status":        order.Status,
		"event_type":    "order_created",
	}

	eventData, _ := json.Marshal(event)
	if err := redisClient.Publish(context.Background(), "orders:events", string(eventData)).Err(); err != nil {
		log.Println("Warning: Failed to publish event:", err)
	}

	respondSuccess(w, http.StatusCreated, "Order created successfully", &order)
}

func listOrdersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, customer_name, product, quantity, price, status, created_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.Product, &order.Quantity, &order.Price, &order.Status, &order.CreatedAt); err != nil {
			log.Println("Error scanning order:", err)
			continue
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    orders,
	})
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var order Order
	query := "SELECT id, customer_name, product, quantity, price, status, created_at FROM orders WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&order.ID, &order.CustomerName, &order.Product, &order.Quantity, &order.Price, &order.Status, &order.CreatedAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Order not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch order")
		return
	}

	respondSuccess(w, http.StatusOK, "Order retrieved", &order)
}

func cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Get order
	var order Order
	query := "SELECT id, customer_name, product, quantity, price, status, created_at FROM orders WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&order.ID, &order.CustomerName, &order.Product, &order.Quantity, &order.Price, &order.Status, &order.CreatedAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Order not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch order")
		return
	}

	// Validate status
	if order.Status == "cancelled" {
		respondError(w, http.StatusBadRequest, "Order already cancelled")
		return
	}
	if order.Status == "shipped" {
		respondError(w, http.StatusBadRequest, "Cannot cancel shipped order")
		return
	}

	// Update status
	_, err = db.Exec("UPDATE orders SET status = $1 WHERE id = $2", "cancelled", id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	order.Status = "cancelled"

	// Publish event
	event := map[string]interface{}{
		"id":            order.ID,
		"customer_name": order.CustomerName,
		"product":       order.Product,
		"quantity":      order.Quantity,
		"price":         order.Price,
		"status":        order.Status,
		"event_type":    "order_cancelled",
	}

	eventData, _ := json.Marshal(event)
	if err := redisClient.Publish(context.Background(), "orders:events", string(eventData)).Err(); err != nil {
		log.Println("Warning: Failed to publish event:", err)
	}

	respondSuccess(w, http.StatusOK, "Order cancelled", &order)
}

func respondSuccess(w http.ResponseWriter, statusCode int, message string, data *Order) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(OrderResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(OrderResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
