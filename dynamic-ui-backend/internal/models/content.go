package models

import "time"

type Category struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	SearchText   string    `json:"search_text"`
	ImageURL     string    `json:"image_url"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    *int      `json:"created_by,omitempty"`
	UpdatedBy    *int      `json:"updated_by,omitempty"`
}

type Brand struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	SearchText   string    `json:"search_text"`
	ImageURL     string    `json:"image_url"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    *int      `json:"created_by,omitempty"`
	UpdatedBy    *int      `json:"updated_by,omitempty"`
}

type CreateCategoryRequest struct {
	Name         string `json:"name"`
	SearchText   string `json:"search_text"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`
	SearchText   *string `json:"search_text,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type CreateBrandRequest struct {
	Name         string `json:"name"`
	SearchText   string `json:"search_text"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

type UpdateBrandRequest struct {
	Name         *string `json:"name,omitempty"`
	SearchText   *string `json:"search_text,omitempty"`
	ImageURL     *string `json:"image_url,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}
