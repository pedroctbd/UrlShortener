package entities

type ShortenedUrl struct {
	OriginalUrl  string `json:"original"`
	ShortenedUrl string `json:"short"`
	CreationDate string `json:"created_at"`
}
