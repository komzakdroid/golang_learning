package handlers

import (
	"dynamic-ui-backend/internal/models"
	"dynamic-ui-backend/internal/services"
	"dynamic-ui-backend/pkg/logger"
	"encoding/json"
	"net/http"
	"time"

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
