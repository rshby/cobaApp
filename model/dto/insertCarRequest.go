package dto

type InsertCarRequest struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0.000"`
	ReleaseDate string  `json:"release_date" validate:"required"`
}
