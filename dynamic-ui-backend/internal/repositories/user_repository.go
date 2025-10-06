package repositories

import (
	"dynamic-ui-backend/internal/database"
	"dynamic-ui-backend/internal/models"
	"fmt"
	"time"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
        SELECT id, username, email, password_hash, role, is_active, created_at, updated_at
        FROM users WHERE username = $1 AND is_active = true
    `, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *UserRepository) CreateSession(userID int, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
        INSERT INTO sessions (user_id, token, expires_at)
        VALUES ($1, $2, $3)
    `, userID, token, expiresAt)
	return err
}

func (r *UserRepository) GetSession(token string) (*models.Session, error) {
	session := &models.Session{}
	err := r.db.QueryRow(`
        SELECT id, user_id, token, expires_at, created_at
        FROM sessions WHERE token = $1 AND expires_at > NOW()
    `, token).Scan(
		&session.ID, &session.UserID, &session.Token,
		&session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	return session, nil
}

func (r *UserRepository) DeleteSession(token string) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}
