package handlers

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"shorturl.com/entities"
	"shorturl.com/utils"
)

func generateShortCode(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// CreateURL godoc
// @Summary      Create a new shortened URL
// @Description  Takes an original URL and returns a shortened one
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        input  body      entities.CreateURLInput  true  "URL creation input"
// @Success      201    {object}  entities.URL
// @Router       /url [post]
func CreateUrl(db *sql.DB) http.HandlerFunc {

	query := `
		INSERT INTO urls (
			id, short_code, original_url, user_id, created_at, expires_at, access_count, last_accessed_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, NOW(), $4, 0, NULL
		)
		RETURNING id
	`
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input, err := utils.Decode[entities.CreateURLInput](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newShortCode := input.ShortCode
		//if code is empty or less than 8, generate random code
		if len(input.ShortCode) <= 7 || input.ShortCode == "" {

			for {
				newShortCode = generateShortCode(8)
				var exists bool

				err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)", newShortCode).Scan(&exists)

				if err != nil {

					http.Error(w, "Error validating short code", http.StatusInternalServerError)
					return
				}

				if !exists {

					break
				}
			}

		}

		// Insert into database
		var newUrlID uuid.UUID
		err = db.QueryRowContext(ctx, query,
			newShortCode,
			input.OriginalURL,
			input.UserId,
			input.ExpiresAt,
		).Scan(&newUrlID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"id": newUrlID,
		}
		if err := utils.Encode(w, r, http.StatusCreated, resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ListExistingUrls godoc
// @Summary      List all shortened URLs
// @Description  Returns a list of all existing shortened URLs
// @Tags         urls
// @Accept       json
// @Produce      json
// @Success      200  {array}   entities.URL
// @Router       /url/list [get]
func ListExistingUrls(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT * FROM urls")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		urls := []entities.URL{}

		for rows.Next() {
			var url entities.URL

			err := rows.Scan(
				&url.ID,
				&url.ShortCode,
				&url.OriginalURL,
				&url.UserID,
				&url.CreatedAt,
				&url.ExpiresAt,
				&url.AccessCount,
				&url.LastAccessedAt,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			urls = append(urls, url)
		}

		if err := utils.Encode(w, r, http.StatusOK, urls); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

// RedirectUrl godoc
// @Summary      Redirect to original URL
// @Description  Redirects user to the original URL based on the shortened one
// @Tags         urls
// @Produce      plain
// @Success      301  "Redirects to original URL"
// @Router       /{code} [get]
func RedirectUrl(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := chi.URLParam(r, "code")

		var originalURL string
		query := `SELECT original_url FROM urls WHERE short_code = $1`
		err := db.QueryRowContext(ctx, query, code).Scan(&originalURL)

		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//updates access_count
		go func() {

			_, err := db.ExecContext(ctx, `UPDATE urls SET access_count = access_count + 1 WHERE short_code = $1`, code)

			if err != nil {
				log.Printf("Error updating access count: %v", err)
			}
		}()

		//postgres to get original url from db
		http.Redirect(w, r, originalURL, http.StatusPermanentRedirect)

	}
}
