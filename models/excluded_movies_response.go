package models

type ExcludedMoviesResponse struct {
	Username         string   `json:"username"`
	ExcludedMovieIds []string `json:"excludedMovieIds"`
}
