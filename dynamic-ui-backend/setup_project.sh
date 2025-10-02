#!/bin/bash

echo "🚀 Creating all files..."

# Create all directories
mkdir -p cmd/server
mkdir -p internal/api/handlers
mkdir -p internal/api/middleware
mkdir -p internal/models
mkdir -p internal/services
mkdir -p pkg/logger
mkdir -p schemas/v1

# FILE 1: cmd/server/main.go
cat > cmd/server/main.go << 'MAINEOF'
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

	uiService := services.NewUIService()
	appLogger.Info("✅ UI Service initialized")

	router := api.SetupRoutes(uiService, appLogger)

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
		appLogger.Info(fmt.Sprintf("📱 Test: http://%s/api/v1/ui/survey", addr))
		
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
MAINEOF

# FILE 2: pkg/logger/logger.go
cat > pkg/logger/logger.go << 'LOGGEREOF'
package logger

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() (*Logger, error) {
	env := os.Getenv("ENV")
	var config zap.Config
	
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(logLevel)); err == nil {
			config.Level = zap.NewAtomicLevelAt(level)
		}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger.Sugar()}, nil
}
LOGGEREOF

# FILE 3: internal/models/ui_schema.go
cat > internal/models/ui_schema.go << 'MODELSEOF'
package models

import "time"

type UISchemaResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Message  string      `json:"message,omitempty"`
	Version  string      `json:"version"`
	CachedAt time.Time   `json:"cached_at"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

type VersionInfo struct {
	Success          bool      `json:"success"`
	AppVersion       string    `json:"app_version"`
	MinVersion       string    `json:"min_version"`
	ForceUpdate      bool      `json:"force_update"`
	AvailableScreens []string  `json:"available_screens"`
	UpdatedAt        time.Time `json:"updated_at"`
}
MODELSEOF

# FILE 4: internal/services/ui_service.go
cat > internal/services/ui_service.go << 'SERVICEEOF'
package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
	"github.com/patrickmn/go-cache"
)

type UIService struct {
	cache      *cache.Cache
	schemaPath string
}

func NewUIService() *UIService {
	schemaPath := os.Getenv("SCHEMA_BASE_PATH")
	if schemaPath == "" {
		schemaPath = "./schemas"
	}
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &UIService{cache: c, schemaPath: schemaPath}
}

func (s *UIService) GetScreenSchema(screenName, version string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("schema_%s_%s", screenName, version)
	if cached, found := s.cache.Get(cacheKey); found {
		return cached.(map[string]interface{}), nil
	}

	filePath := filepath.Join(s.schemaPath, version, fmt.Sprintf("%s.json", screenName))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("schema not found: %w", err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	s.cache.Set(cacheKey, schema, cache.DefaultExpiration)
	return schema, nil
}

func (s *UIService) GetAvailableScreens(version string) ([]string, error) {
	dirPath := filepath.Join(s.schemaPath, version)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	screens := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			screenName := file.Name()[:len(file.Name())-5]
			screens = append(screens, screenName)
		}
	}
	return screens, nil
}

func (s *UIService) ClearCache() {
	s.cache.Flush()
}
SERVICEEOF

# FILE 5: internal/api/middleware/cors.go
cat > internal/api/middleware/cors.go << 'CORSEOF'
package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if origins == "" {
			origins = "*"
		}

		methods := os.Getenv("CORS_ALLOWED_METHODS")
		if methods == "" {
			methods = "GET,POST,PUT,DELETE,OPTIONS"
		}

		headers := os.Getenv("CORS_ALLOWED_HEADERS")
		if headers == "" {
			headers = "Content-Type,Authorization"
		}

		if origins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origin := r.Header.Get("Origin")
			allowedOrigins := strings.Split(origins, ",")
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", methods)
		w.Header().Set("Access-Control-Allow-Headers", headers)
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
CORSEOF

# FILE 6: internal/api/middleware/logger.go
cat > internal/api/middleware/logger.go << 'LOGMIDEOF'
package middleware

import (
	"dynamic-ui-backend/pkg/logger"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func Logger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(rw, r)
			duration := time.Since(start)

			log.Infow("HTTP Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", duration.String(),
				"size", rw.size,
			)
		})
	}
}
LOGMIDEOF

# FILE 7: internal/api/handlers/health.go
cat > internal/api/handlers/health.go << 'HEALTHEOF'
package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startTime: time.Now()}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":  "healthy",
		"message": "Server is running",
		"time":    time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) DetailedHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	response := map[string]interface{}{
		"status": "healthy",
		"uptime": uptime.String(),
		"system": map[string]interface{}{
			"goroutines":   runtime.NumGoroutine(),
			"memory_alloc": memStats.Alloc / 1024 / 1024,
			"memory_total": memStats.TotalAlloc / 1024 / 1024,
		},
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
HEALTHEOF

# FILE 8: internal/api/handlers/ui_handler.go
cat > internal/api/handlers/ui_handler.go << 'UIHANDLEREOF'
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"dynamic-ui-backend/internal/models"
	"dynamic-ui-backend/internal/services"
	"dynamic-ui-backend/pkg/logger"
	"github.com/gorilla/mux"
)

type UIHandler struct {
	uiService *services.UIService
	logger    *logger.Logger
}

func NewUIHandler(uiService *services.UIService, log *logger.Logger) *UIHandler {
	return &UIHandler{uiService: uiService, logger: log}
}

func (h *UIHandler) GetScreen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	screenName := vars["screen"]
	version := r.URL.Query().Get("version")
	if version == "" {
		version = "v1"
	}

	schema, err := h.uiService.GetScreenSchema(screenName, version)
	if err != nil {
		h.logger.Errorw("Failed to get schema", "screen", screenName, "error", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Schema not found",
			Code:    "SCHEMA_NOT_FOUND",
		})
		return
	}

	response := models.UISchemaResponse{
		Success:  true,
		Data:     schema,
		Version:  version,
		CachedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UIHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")
	if version == "" {
		version = "v1"
	}

	screens, _ := h.uiService.GetAvailableScreens(version)

	response := models.VersionInfo{
		Success:          true,
		AppVersion:       "1.0.0",
		MinVersion:       "1.0.0",
		ForceUpdate:      false,
		AvailableScreens: screens,
		UpdatedAt:        time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UIHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	h.uiService.ClearCache()
	response := map[string]interface{}{"success": true, "message": "Cache cleared"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	h.logger.Info("Cache cleared")
}

func (h *UIHandler) ListScreens(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")
	if version == "" {
		version = "v1"
	}

	screens, err := h.uiService.GetAvailableScreens(version)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Success: false, Error: "Failed to get screens"})
		return
	}

	response := map[string]interface{}{"success": true, "screens": screens, "version": version}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
UIHANDLEREOF

# FILE 9: internal/api/routes.go
cat > internal/api/routes.go << 'ROUTESEOF'
package api

import (
	"dynamic-ui-backend/internal/api/handlers"
	"dynamic-ui-backend/internal/api/middleware"
	"dynamic-ui-backend/internal/services"
	"dynamic-ui-backend/pkg/logger"
	"github.com/gorilla/mux"
)

func SetupRoutes(uiService *services.UIService, log *logger.Logger) *mux.Router {
	router := mux.NewRouter()

	uiHandler := handlers.NewUIHandler(uiService, log)
	healthHandler := handlers.NewHealthHandler()

	router.Use(middleware.CORS)
	router.Use(middleware.Logger(log))

	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	router.HandleFunc("/health/detailed", healthHandler.DetailedHealth).Methods("GET")

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/ui/{screen}", uiHandler.GetScreen).Methods("GET")
	api.HandleFunc("/ui/version", uiHandler.GetVersion).Methods("GET")
	api.HandleFunc("/ui/screens", uiHandler.ListScreens).Methods("GET")
	api.HandleFunc("/admin/cache/clear", uiHandler.ClearCache).Methods("POST")

	return router
}
ROUTESEOF

echo "✅ All Go files created!"
echo ""
echo "Now creating JSON schemas..."

# Create survey.json schema
cat > schemas/v1/survey.json << 'SURVEYJSON'
{
  "screen_id": "survey_v1",
  "version": "1.0.0",
  "title": "Bosh Sahifa",
  "background_color": "#F5F7FA",
  "app_bar": {
    "title": "Dynamic UI Demo",
    "background_color": "#2563EB",
    "text_color": "#FFFFFF",
    "elevation": 2
  },
  "widgets": [
    {
      "id": "header_card",
      "type": "container",
      "padding": {"top": 20, "bottom": 20, "left": 16, "right": 16},
      "margin": {"all": 16},
      "decoration": {
        "background_color": "#FFFFFF",
        "border_radius": 16,
        "shadow": {"color": "#00000020", "blur_radius": 10, "offset": {"x": 0, "y": 4}}
      },
      "children": [
        {
          "type": "text",
          "content": "Xush kelibsiz! 👋",
          "style": {"font_size": 28, "font_weight": "bold", "color": "#1F2937"}
        },
        {"type": "sized_box", "height": 8},
        {
          "type": "text",
          "content": "Server-driven UI namunasi",
          "style": {"font_size": 16, "color": "#6B7280"}
        }
      ]
    }
  ]
}
SURVEYJSON

echo "✅ JSON schemas created!"
echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. go mod tidy"
echo "2. go run cmd/server/main.go"

