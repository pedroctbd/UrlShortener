package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"shorturl.com/entities"
	"shorturl.com/utils"
)

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

		// Insert into database
		var newUrlID uuid.UUID
		err = db.QueryRowContext(ctx, query,
			input.ShortCode,
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
// @Router       /urls [get]
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
// @Router       /redirect [get]
func RedirectUrl() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//postgres to get original url from db
		http.Redirect(w, r, "https://youtube.com", http.StatusPermanentRedirect)

	}
}
