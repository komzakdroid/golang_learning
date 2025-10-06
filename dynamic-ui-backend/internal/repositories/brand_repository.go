package repositories

import (
	"dynamic-ui-backend/internal/database"
	"dynamic-ui-backend/internal/models"
	"fmt"
)

type BrandRepository struct {
	db *database.DB
}

func NewBrandRepository(db *database.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

func (r *BrandRepository) GetAll() ([]models.Brand, error) {
	rows, err := r.db.Query(`
        SELECT id, name, search_text, image_url, display_order, is_active,
               created_at, updated_at, created_by, updated_by
        FROM brands
        WHERE is_active = true
        ORDER BY display_order ASC, id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brands := make([]models.Brand, 0)
	for rows.Next() {
		var brand models.Brand
		err := rows.Scan(
			&brand.ID, &brand.Name, &brand.SearchText, &brand.ImageURL,
			&brand.DisplayOrder, &brand.IsActive, &brand.CreatedAt, &brand.UpdatedAt,
			&brand.CreatedBy, &brand.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}
	return brands, nil
}

func (r *BrandRepository) GetByID(id int) (*models.Brand, error) {
	brand := &models.Brand{}
	err := r.db.QueryRow(`
        SELECT id, name, search_text, image_url, display_order, is_active,
               created_at, updated_at, created_by, updated_by
        FROM brands WHERE id = $1
    `, id).Scan(
		&brand.ID, &brand.Name, &brand.SearchText, &brand.ImageURL,
		&brand.DisplayOrder, &brand.IsActive, &brand.CreatedAt, &brand.UpdatedAt,
		&brand.CreatedBy, &brand.UpdatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}
	return brand, nil
}

func (r *BrandRepository) Create(req *models.CreateBrandRequest, userID int) (*models.Brand, error) {
	brand := &models.Brand{}
	err := r.db.QueryRow(`
        INSERT INTO brands (name, search_text, image_url, display_order, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $5)
        RETURNING id, name, search_text, image_url, display_order, is_active, created_at, updated_at
    `, req.Name, req.SearchText, req.ImageURL, req.DisplayOrder, userID).Scan(
		&brand.ID, &brand.Name, &brand.SearchText, &brand.ImageURL,
		&brand.DisplayOrder, &brand.IsActive, &brand.CreatedAt, &brand.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return brand, nil
}

func (r *BrandRepository) Update(id int, req *models.UpdateBrandRequest, userID int) (*models.Brand, error) {
	query := `UPDATE brands SET updated_by = $1, updated_at = NOW()`
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

	brand := &models.Brand{}
	err := r.db.QueryRow(query, args...).Scan(
		&brand.ID, &brand.Name, &brand.SearchText, &brand.ImageURL,
		&brand.DisplayOrder, &brand.IsActive, &brand.CreatedAt, &brand.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return brand, nil
}

func (r *BrandRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM brands WHERE id = $1`, id)
	return err
}
