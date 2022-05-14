package models

type LikedMovieRequest struct {
	Username string `json:"username" validate:"required"`
	MovieId  string `json:"movieId" validate:"required"`
}
