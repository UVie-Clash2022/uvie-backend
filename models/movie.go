package models

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	PosterUrl   string `json:"poster_path"`
}
