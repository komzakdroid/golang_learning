package api

import (
	"dynamic-ui-backend/internal/api/handlers"
	"dynamic-ui-backend/internal/api/middleware"
	"dynamic-ui-backend/internal/database"
	"dynamic-ui-backend/internal/repositories"
	"dynamic-ui-backend/internal/services"
	"dynamic-ui-backend/pkg/logger"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *database.DB, uiService *services.UIService, log *logger.Logger) *mux.Router {
	router := mux.NewRouter()

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	brandRepo := repositories.NewBrandRepository(db)

	// Handlers
	uiHandler := handlers.NewUIHandler(uiService, log)
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(userRepo, log)
	adminHandler := handlers.NewAdminHandler(categoryRepo, brandRepo, log)
	uploadHandler := handlers.NewUploadHandler(log)

	// Global middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logger(log))

	// Health
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Static files - uploads directory
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	router.PathPrefix("/uploads/").Handler(
		http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))),
	).Methods("GET")

	// API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// Auth endpoints (public)
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// UI Schema (public)
	api.HandleFunc("/ui", uiHandler.GetScreen).Methods("GET")
	api.HandleFunc("/ui/version", uiHandler.GetVersion).Methods("GET")

	// Content endpoints (public - for mobile app)
	api.HandleFunc("/content/categories", adminHandler.GetAllCategories).Methods("GET")
	api.HandleFunc("/content/brands", adminHandler.GetAllBrands).Methods("GET")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	// Admin only routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminOnly)

	// File upload
	admin.HandleFunc("/upload", uploadHandler.UploadImage).Methods("POST")

	// Categories management
	admin.HandleFunc("/categories", adminHandler.GetAllCategories).Methods("GET")
	admin.HandleFunc("/categories", adminHandler.CreateCategory).Methods("POST")
	admin.HandleFunc("/categories/{id}", adminHandler.GetCategory).Methods("GET")
	admin.HandleFunc("/categories/{id}", adminHandler.UpdateCategory).Methods("PUT")
	admin.HandleFunc("/categories/{id}", adminHandler.DeleteCategory).Methods("DELETE")

	// Brands management
	admin.HandleFunc("/brands", adminHandler.GetAllBrands).Methods("GET")
	admin.HandleFunc("/brands", adminHandler.CreateBrand).Methods("POST")
	admin.HandleFunc("/brands/{id}", adminHandler.UpdateBrand).Methods("PUT")
	admin.HandleFunc("/brands/{id}", adminHandler.DeleteBrand).Methods("DELETE")

	admin.HandleFunc("/cache/clear", uiHandler.ClearCache).Methods("POST")

	return router
}
