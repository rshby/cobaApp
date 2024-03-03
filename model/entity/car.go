package entity

import (
	"database/sql"
)

type Car struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	ReleaseDate *sql.NullTime `json:"release_date"`
}
