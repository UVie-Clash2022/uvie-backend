package models

type RecommendMoviesResponse struct {
	Page         int     `json:"page"`
	Movies       []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}
