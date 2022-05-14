package models

type UserMoviesLiked struct {
	Username        string  `json:"username,omitempty" bson:"username"`
	CurrentFavorite Movie   `json:"currentFavorite" bson:"currentFavorite"`
	Liked           []Movie `json:"liked" bson:"liked"`
}
