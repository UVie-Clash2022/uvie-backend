package models

type ExcludeMovieRequest struct {
	Username string `json:"username" validate:"required"`
	MovieId  string `json:"movieId" validate:"required"`
}
