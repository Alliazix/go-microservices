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

// Notification represents a notification record
type Notification struct {
	ID        string     `json:"id"`
	OrderID   string     `json:"order_id"`
	Message   string     `json:"message"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
}

type NotificationResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *Notification `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
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

	// Start listening for order events in a goroutine
	go listenForOrderEvents()

	// Setup HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/notifications", listNotificationsHandler).Methods("GET")
	router.HandleFunc("/api/notifications/{id}", getNotificationHandler).Methods("GET")

	port := getEnv("PORT", "8081")
	log.Printf("✓ Notification Service starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id UUID NOT NULL,
		message TEXT NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		sent_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_notifications_order_id ON notifications(order_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
	`

	_, err := db.Exec(query)
	return err
}

// listenForOrderEvents subscribes to Redis and processes order events
func listenForOrderEvents() {
	pubsub := redisClient.Subscribe(context.Background(), "orders:events")
	defer pubsub.Close()

	ch := pubsub.Channel()
	log.Println("✓ Listening for order events on Redis...")

	for msg := range ch {
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Println("Error parsing event:", err)
			continue
		}

		// Process the event
		handleOrderEvent(event)
	}
}

func handleOrderEvent(event map[string]interface{}) {
	eventType, _ := event["event_type"].(string)
	orderID, _ := event["id"].(string)
	customerName, _ := event["customer_name"].(string)
	product, _ := event["product"].(string)
	quantity, _ := event["quantity"].(float64)
	price, _ := event["price"].(float64)

	var message string
	switch eventType {
	case "order_created":
		message = fmt.Sprintf("Order created: %s x%.0f ($%.2f) for %s", product, quantity, price, customerName)
	case "order_cancelled":
		message = fmt.Sprintf("Order cancelled for %s", customerName)
	default:
		log.Printf("Unknown event type: %s\n", eventType)
		return
	}

	// Create notification in database
	notification := Notification{
		OrderID: orderID,
		Message: message,
		Status:  "pending",
	}

	query := `
		INSERT INTO notifications (order_id, message, status)
		VALUES ($1, $2, $3)
		RETURNING id, order_id, message, status, created_at, sent_at
	`

	var sentAt *time.Time
	err := db.QueryRow(query, notification.OrderID, notification.Message, "pending").
		Scan(&notification.ID, &notification.OrderID, &notification.Message, &notification.Status, &notification.CreatedAt, &sentAt)

	if err != nil {
		log.Println("Error creating notification:", err)
		return
	}

	// Simulate sending and mark as sent
	time.Sleep(100 * time.Millisecond)

	now := time.Now()
	_, err = db.Exec("UPDATE notifications SET status = $1, sent_at = $2 WHERE id = $3", "sent", now, notification.ID)
	if err != nil {
		log.Println("Error updating notification:", err)
		return
	}

	log.Printf("✓ Notification created: %s\n", message)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func listNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, order_id, message, status, created_at, sent_at FROM notifications ORDER BY created_at DESC")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notif Notification
		var sentAt *time.Time
		if err := rows.Scan(&notif.ID, &notif.OrderID, &notif.Message, &notif.Status, &notif.CreatedAt, &sentAt); err != nil {
			log.Println("Error scanning notification:", err)
			continue
		}
		notif.SentAt = sentAt
		notifications = append(notifications, notif)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    notifications,
	})
}

func getNotificationHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var notif Notification
	var sentAt *time.Time
	query := "SELECT id, order_id, message, status, created_at, sent_at FROM notifications WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&notif.ID, &notif.OrderID, &notif.Message, &notif.Status, &notif.CreatedAt, &sentAt)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Notification not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch notification")
		return
	}

	notif.SentAt = sentAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NotificationResponse{
		Success: true,
		Message: "Notification retrieved",
		Data:    &notif,
	})
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(NotificationResponse{
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
