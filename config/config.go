package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
)

func InitDB() *sql.DB {
	connStr := "host=db user=postgres password=secret dbname=messaging sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Auto-migrate the messages table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		to TEXT NOT NULL,
		content TEXT NOT NULL CHECK (char_length(content) <= 160),
		is_sent BOOLEAN DEFAULT FALSE,
		message_id TEXT,
		sent_at TIMESTAMP
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	return db
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}
