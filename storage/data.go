package storage

import (
	"time"

	"shorturl.com/entities"
)

var CurrentUrls = []entities.ShortenedUrl{
	{
		ShortenedUrl: "teste.com",
		OriginalUrl:  "testeoriginal.com",
		CreationDate: time.Now().Format(time.RFC3339),
	},
}
