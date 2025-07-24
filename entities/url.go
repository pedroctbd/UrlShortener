package entities

import (
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ShortCode      string     `db:"short_code" json:"short_code"`
	OriginalURL    string     `db:"original_url" json:"original_url"`
	UserID         *uuid.UUID `db:"user_id" json:"user_id,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	ExpiresAt      *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	AccessCount    int64      `db:"access_count" json:"access_count"`
	LastAccessedAt *time.Time `db:"last_accessed_at" json:"last_accessed_at,omitempty"`
}

type CreateURLInput struct {
	ShortCode   string     `json:"short_code"`
	UserId      string     `json:"user_id"`
	OriginalURL string     `json:"original_url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
