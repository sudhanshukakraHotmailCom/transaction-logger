package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"transaction-logger/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDB(cfg *config.Config) (*Database, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *sql.DB
	var err error
	maxRetries := 10
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Failed to open database: %v. Retrying in %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return &Database{DB: db}, nil
		}

		log.Printf("Failed to ping database (attempt %d/%d): %v. Retrying in %v...",
			i+1, maxRetries, err, retryDelay)
		db.Close()
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}

func (d *Database) Close() error {
	return d.DB.Close()
}

// InitSchema initializes the database schema
func (db *Database) InitSchema() error {
	// Create users table
	_, err := db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`)

	if err != nil {
		return err
	}

	// Create transactions table with foreign key to users
	_, err = db.DB.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			sender_account TEXT NOT NULL,
			receiver_account TEXT NOT NULL,
			amount DECIMAL(15, 2) NOT NULL,
			currency TEXT NOT NULL,
			transaction_type TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id TEXT NOT NULL,
			CONSTRAINT fk_user
				FOREIGN KEY (user_id)
				REFERENCES users(id)
				ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
	`)

	return err
}
