package dto

type InsertCarResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ReleaseDate string  `json:"release_date"`
}
