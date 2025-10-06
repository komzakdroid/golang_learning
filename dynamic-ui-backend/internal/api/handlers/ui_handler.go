package handlers

import (
	"dynamic-ui-backend/internal/models"
	"dynamic-ui-backend/internal/services"
	"dynamic-ui-backend/pkg/logger"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type UIHandler struct {
	uiService *services.UIService
	logger    *logger.Logger
}

func NewUIHandler(uiService *services.UIService, log *logger.Logger) *UIHandler {
	return &UIHandler{uiService: uiService, logger: log}
}

func (h *UIHandler) GetScreen(w http.ResponseWriter, r *http.Request) {
	screenName := r.URL.Query().Get("screen")

	if screenName == "" {
		h.logger.Error("Screen parameter is missing")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Screen parameter is required",
			Code:    "MISSING_PARAMETER",
		})
		return
	}

	if strings.Contains(screenName, "..") || strings.Contains(screenName, "/") || strings.Contains(screenName, "\\") {
		h.logger.Errorw("Invalid screen name", "screen", screenName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Invalid screen name",
			Code:    "INVALID_PARAMETER",
		})
		return
	}

	version := r.URL.Query().Get("version")
	if version == "" {
		version = "v1"
	}

	schema, err := h.uiService.GetScreenSchema(screenName, version)
	if err != nil {
		h.logger.Errorw("Failed to get schema", "screen", screenName, "version", version, "error", err)
		w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Success: false, Error: "Failed to get screens"})
		return
	}

	response := map[string]interface{}{"success": true, "screens": screens, "version": version}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
