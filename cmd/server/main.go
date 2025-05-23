package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"transaction-logger/internal/auth"
	"transaction-logger/internal/config"
	"transaction-logger/internal/database"
	"transaction-logger/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize auth package with config
	auth.Init(cfg)

	// Initialize database schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Initialize handlers
	transactionHandler := handlers.NewTransactionHandler(db.DB)
	authHandler := handlers.NewAuthHandler(db.DB)

	// API router with auth middleware
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(handlers.AuthMiddleware)

	// Auth routes (public)
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Transaction routes (protected by auth middleware)
	apiRouter.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods("GET")
	apiRouter.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")
	apiRouter.HandleFunc("/transactions/generatesample", transactionHandler.GenerateSampleTransactions).Methods("POST")

	// Start server
	port := ":" + cfg.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Printf("JWT Secret: %s", auth.GetJWTSecret())
	log.Fatal(http.ListenAndServe(port, router))
}
