package handlers

import (
	"net/http"

	"shorturl.com/entities"
	"shorturl.com/storage"
	"shorturl.com/utils"
)

// CreateUrl godoc
// @Summary      Create a shortened URL
// @Description  Accepts a JSON payload to create a new shortened URL
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        url  body      entities.ShortenedUrl  true  "URL info"
// @Success      201  {array}   entities.ShortenedUrl
// @Failure      400  {string}  string  "Invalid request payload"
// @Failure      500  {string}  string  "Internal server error"
// @Router       / [post]
func CreateUrl() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		newUrl, err := utils.Decode[entities.ShortenedUrl](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		storage.CurrentUrls = append(storage.CurrentUrls, newUrl)

		if err := utils.Encode(w, r, http.StatusCreated, storage.CurrentUrls); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RedirectUrl godoc
// @Summary      Redirect to original URL
// @Description  Redirects user to the original URL based on the shortened one
// @Tags         urls
// @Produce      plain
// @Success      301  "Redirects to original URL"
// @Router       / [get]
func RedirectUrl() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//postgres to get original url from db

		http.Redirect(w, r, "https://google.com", http.StatusPermanentRedirect)

	}
}
