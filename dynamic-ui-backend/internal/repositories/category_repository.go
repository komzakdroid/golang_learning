package repositories

import (
	"dynamic-ui-backend/internal/database"
	"dynamic-ui-backend/internal/models"
	"fmt"
)

type CategoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query(`
        SELECT id, name, search_text, image_url, display_order, is_active,
               created_at, updated_at, created_by, updated_by
        FROM categories
        WHERE is_active = true
        ORDER BY display_order ASC, id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var cat models.Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.SearchText, &cat.ImageURL,
			&cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
			&cat.CreatedBy, &cat.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	cat := &models.Category{}
	err := r.db.QueryRow(`
        SELECT id, name, search_text, image_url, display_order, is_active,
               created_at, updated_at, created_by, updated_by
        FROM categories WHERE id = $1
    `, id).Scan(
		&cat.ID, &cat.Name, &cat.SearchText, &cat.ImageURL,
		&cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
		&cat.CreatedBy, &cat.UpdatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return cat, nil
}

func (r *CategoryRepository) Create(req *models.CreateCategoryRequest, userID int) (*models.Category, error) {
	cat := &models.Category{}
	err := r.db.QueryRow(`
        INSERT INTO categories (name, search_text, image_url, display_order, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $5)
        RETURNING id, name, search_text, image_url, display_order, is_active, created_at, updated_at
    `, req.Name, req.SearchText, req.ImageURL, req.DisplayOrder, userID).Scan(
		&cat.ID, &cat.Name, &cat.SearchText, &cat.ImageURL,
		&cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *CategoryRepository) Update(id int, req *models.UpdateCategoryRequest, userID int) (*models.Category, error) {
	query := `UPDATE categories SET updated_by = $1, updated_at = NOW()`
	args := []interface{}{userID}
	argPos := 2

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argPos)
		args = append(args, *req.Name)
		argPos++
	}
	if req.SearchText != nil {
		query += fmt.Sprintf(", search_text = $%d", argPos)
		args = append(args, *req.SearchText)
		argPos++
	}
	if req.ImageURL != nil {
		query += fmt.Sprintf(", image_url = $%d", argPos)
		args = append(args, *req.ImageURL)
		argPos++
	}
	if req.DisplayOrder != nil {
		query += fmt.Sprintf(", display_order = $%d", argPos)
		args = append(args, *req.DisplayOrder)
		argPos++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argPos)
		args = append(args, *req.IsActive)
		argPos++
	}

	query += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, search_text, image_url, display_order, is_active, created_at, updated_at", argPos)
	args = append(args, id)

	cat := &models.Category{}
	err := r.db.QueryRow(query, args...).Scan(
		&cat.ID, &cat.Name, &cat.SearchText, &cat.ImageURL,
		&cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *CategoryRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}
