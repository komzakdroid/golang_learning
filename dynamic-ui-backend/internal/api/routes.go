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
