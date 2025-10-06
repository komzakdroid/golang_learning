package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dynamic-ui-backend/internal/api"
	"dynamic-ui-backend/internal/database"
	"dynamic-ui-backend/internal/services"
	"dynamic-ui-backend/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	appLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer appLogger.Sync()

	appLogger.Info("🚀 Starting Dynamic UI Backend Server...")

	// Database
	db, err := database.NewDB()
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}
	defer db.Close()
	appLogger.Info("✅ Database connected")

	// Services
	uiService := services.NewUIService()
	appLogger.Info("✅ UI Service initialized")

	// Routes
	router := api.SetupRoutes(db, uiService, appLogger)

	port := getEnv("SERVER_PORT", "8080")
	host := getEnv("SERVER_HOST", "0.0.0.0")
	addr := fmt.Sprintf("%s:%s", host, port)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		appLogger.Info(fmt.Sprintf("🌐 Server listening on http://%s", addr))
		appLogger.Info(fmt.Sprintf("📱 API: http://%s/api/v1", addr))
		appLogger.Info(fmt.Sprintf("💚 Health: http://%s/health", addr))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal(fmt.Sprintf("Server failed: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("🛑 Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal(fmt.Sprintf("Shutdown error: %v", err))
	}

	appLogger.Info("✅ Stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
