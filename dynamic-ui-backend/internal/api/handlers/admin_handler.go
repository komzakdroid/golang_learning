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
	"strconv"
	"strings"
	"time"

	"dynamic-ui-backend/internal/auth"
	"dynamic-ui-backend/internal/models"
	"dynamic-ui-backend/internal/repositories"
	"dynamic-ui-backend/pkg/logger"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
	categoryRepo *repositories.CategoryRepository
	brandRepo    *repositories.BrandRepository
	logger       *logger.Logger
	uploadDir    string
	baseURL      string
}

func NewAdminHandler(
	categoryRepo *repositories.CategoryRepository,
	brandRepo *repositories.BrandRepository,
	log *logger.Logger,
) *AdminHandler {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Create upload directories
	os.MkdirAll(filepath.Join(uploadDir, "categories"), 0755)
	os.MkdirAll(filepath.Join(uploadDir, "brands"), 0755)

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &AdminHandler{
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		logger:       log,
		uploadDir:    uploadDir,
		baseURL:      baseURL,
	}
}

// Categories with Image Upload
func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*auth.Claims)

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.respondError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get form fields
	name := r.FormValue("name")
	searchText := r.FormValue("search_text")
	displayOrder, _ := strconv.Atoi(r.FormValue("display_order"))

	if name == "" || searchText == "" {
		h.respondError(w, "Name and search_text are required", http.StatusBadRequest)
		return
	}

	// Handle file upload
	file, header, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		// Validate image
		if !h.isValidImage(header.Filename) {
			h.respondError(w, "Invalid image format", http.StatusBadRequest)
			return
		}

		// Save image
		imageURL, err = h.saveImage(file, header.Filename, "categories")
		if err != nil {
			h.logger.Errorw("Failed to save image", "error", err)
			h.respondError(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	} else {
		h.respondError(w, "Image file is required", http.StatusBadRequest)
		return
	}

	// Create category
	req := &models.CreateCategoryRequest{
		Name:         name,
		SearchText:   searchText,
		ImageURL:     imageURL,
		DisplayOrder: displayOrder,
	}

	category, err := h.categoryRepo.Create(req, claims.UserID)
	if err != nil {
		h.logger.Errorw("Failed to create category", "error", err)
		h.respondError(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	h.respondSuccess(w, category)
	h.logger.Infow("Category created", "id", category.ID, "name", category.Name, "by", claims.Username)
}

func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*auth.Claims)
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.respondError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.respondError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	req := &models.UpdateCategoryRequest{}

	// Get form fields
	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}
	if searchText := r.FormValue("search_text"); searchText != "" {
		req.SearchText = &searchText
	}
	if orderStr := r.FormValue("display_order"); orderStr != "" {
		order, _ := strconv.Atoi(orderStr)
		req.DisplayOrder = &order
	}
	if activeStr := r.FormValue("is_active"); activeStr != "" {
		active := activeStr == "true"
		req.IsActive = &active
	}

	// Handle file upload if provided
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		if !h.isValidImage(header.Filename) {
			h.respondError(w, "Invalid image format", http.StatusBadRequest)
			return
		}

		imageURL, err := h.saveImage(file, header.Filename, "categories")
		if err != nil {
			h.logger.Errorw("Failed to save image", "error", err)
			h.respondError(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		req.ImageURL = &imageURL
	}

	category, err := h.categoryRepo.Update(id, req, claims.UserID)
	if err != nil {
		h.logger.Errorw("Failed to update category", "error", err)
		h.respondError(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	h.respondSuccess(w, category)
	h.logger.Infow("Category updated", "id", id, "by", claims.Username)
}

// Brands with Image Upload
func (h *AdminHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*auth.Claims)

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.respondError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get form fields
	name := r.FormValue("name")
	searchText := r.FormValue("search_text")
	displayOrder, _ := strconv.Atoi(r.FormValue("display_order"))

	if name == "" || searchText == "" {
		h.respondError(w, "Name and search_text are required", http.StatusBadRequest)
		return
	}

	// Handle file upload
	file, header, err := r.FormFile("image")
	var imageURL string

	if err == nil {
		defer file.Close()

		if !h.isValidImage(header.Filename) {
			h.respondError(w, "Invalid image format", http.StatusBadRequest)
			return
		}

		imageURL, err = h.saveImage(file, header.Filename, "brands")
		if err != nil {
			h.logger.Errorw("Failed to save image", "error", err)
			h.respondError(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	} else {
		h.respondError(w, "Image file is required", http.StatusBadRequest)
		return
	}

	// Create brand
	req := &models.CreateBrandRequest{
		Name:         name,
		SearchText:   searchText,
		ImageURL:     imageURL,
		DisplayOrder: displayOrder,
	}

	brand, err := h.brandRepo.Create(req, claims.UserID)
	if err != nil {
		h.respondError(w, "Failed to create brand", http.StatusInternalServerError)
		return
	}

	h.respondSuccess(w, brand)
	h.logger.Infow("Brand created", "id", brand.ID, "name", brand.Name)
}

func (h *AdminHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*auth.Claims)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.respondError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	req := &models.UpdateBrandRequest{}

	if name := r.FormValue("name"); name != "" {
		req.Name = &name
	}
	if searchText := r.FormValue("search_text"); searchText != "" {
		req.SearchText = &searchText
	}
	if orderStr := r.FormValue("display_order"); orderStr != "" {
		order, _ := strconv.Atoi(orderStr)
		req.DisplayOrder = &order
	}
	if activeStr := r.FormValue("is_active"); activeStr != "" {
		active := activeStr == "true"
		req.IsActive = &active
	}

	// Handle file upload if provided
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		if !h.isValidImage(header.Filename) {
			h.respondError(w, "Invalid image format", http.StatusBadRequest)
			return
		}

		imageURL, err := h.saveImage(file, header.Filename, "brands")
		if err != nil {
			h.respondError(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		req.ImageURL = &imageURL
	}

	brand, err := h.brandRepo.Update(id, req, claims.UserID)
	if err != nil {
		h.respondError(w, "Failed to update brand", http.StatusInternalServerError)
		return
	}

	h.respondSuccess(w, brand)
}

// Read operations (no changes needed)
func (h *AdminHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryRepo.GetAll()
	if err != nil {
		h.respondError(w, "Failed to get categories", http.StatusInternalServerError)
		return
	}
	h.respondSuccess(w, categories)
}

func (h *AdminHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	category, err := h.categoryRepo.GetByID(id)
	if err != nil {
		h.respondError(w, "Category not found", http.StatusNotFound)
		return
	}
	h.respondSuccess(w, category)
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.categoryRepo.Delete(id); err != nil {
		h.respondError(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}
	h.respondSuccess(w, map[string]string{"message": "Category deleted"})
}

func (h *AdminHandler) GetAllBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.brandRepo.GetAll()
	if err != nil {
		h.respondError(w, "Failed to get brands", http.StatusInternalServerError)
		return
	}
	h.respondSuccess(w, brands)
}

func (h *AdminHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.brandRepo.Delete(id); err != nil {
		h.respondError(w, "Failed to delete brand", http.StatusInternalServerError)
		return
	}
	h.respondSuccess(w, map[string]string{"message": "Brand deleted"})
}

// Helper methods
func (h *AdminHandler) saveImage(file io.Reader, filename, subdir string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	timestamp := time.Now().Unix()
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s-%d", filename, timestamp)))
	hash := hex.EncodeToString(hasher.Sum(nil))[:12]
	newFilename := fmt.Sprintf("%d_%s%s", timestamp, hash, ext)

	// Create file path
	filePath := filepath.Join(h.uploadDir, subdir, newFilename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// Return URL
	return fmt.Sprintf("%s/uploads/%s/%s", h.baseURL, subdir, newFilename), nil
}

func (h *AdminHandler) isValidImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, valid := range validExts {
		if ext == valid {
			return true
		}
	}
	return false
}

func (h *AdminHandler) respondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func (h *AdminHandler) respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Success: false,
		Error:   message,
	})
}
