package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dynamic-ui-backend/internal/auth"
	"dynamic-ui-backend/pkg/logger"
)

type UploadHandler struct {
	uploadDir string
	baseURL   string
	logger    *logger.Logger
}

func NewUploadHandler(log *logger.Logger) *UploadHandler {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Create uploads directory structure
	dirs := []string{
		filepath.Join(uploadDir, "categories"),
		filepath.Join(uploadDir, "brands"),
		filepath.Join(uploadDir, "banners"),
	}
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &UploadHandler{
		uploadDir: uploadDir,
		baseURL:   baseURL,
		logger:    log,
	}
}

type UploadResponse struct {
	Success  bool   `json:"success"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*auth.Claims)

	// Parse multipart form (10MB max)
	maxSize := int64(10 << 20) // 10MB
	if err := r.ParseMultipartForm(maxSize); err != nil {
		h.respondError(w, "File too large or invalid", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.respondError(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get upload type (category, brand, banner)
	uploadType := r.FormValue("type")
	if uploadType == "" {
		uploadType = "general"
	}

	// Validate file type
	if !h.isValidImage(header.Filename) {
		h.respondError(w, "Invalid file type. Only jpg, jpeg, png, gif, webp allowed", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	filename := h.generateFilename(header.Filename)

	// Create subdirectory path
	var subdir string
	switch uploadType {
	case "category":
		subdir = "categories"
	case "brand":
		subdir = "brands"
	case "banner":
		subdir = "banners"
	default:
		subdir = "general"
	}

	filepath := filepath.Join(h.uploadDir, subdir, filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		h.logger.Errorw("Failed to create file", "error", err)
		h.respondError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		h.logger.Errorw("Failed to write file", "error", err)
		h.respondError(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Generate URL
	url := fmt.Sprintf("%s/uploads/%s/%s", h.baseURL, subdir, filename)

	response := UploadResponse{
		Success:  true,
		URL:      url,
		Filename: filename,
		Size:     size,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	h.logger.Infow("File uploaded",
		"filename", filename,
		"size", size,
		"type", uploadType,
		"by", claims.Username,
	)
}

func (h *UploadHandler) isValidImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, valid := range validExts {
		if ext == valid {
			return true
		}
	}
	return false
}

func (h *UploadHandler) generateFilename(original string) string {
	ext := filepath.Ext(original)
	timestamp := time.Now().Unix()

	// Create hash from original filename + timestamp
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%d", original, timestamp)))
	hash := hex.EncodeToString(hasher.Sum(nil))[:12]

	return fmt.Sprintf("%d_%s%s", timestamp, hash, ext)
}

func (h *UploadHandler) respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
