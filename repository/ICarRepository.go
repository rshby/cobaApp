package repository

import (
	"cobaApp/model/entity"
	"context"
	"database/sql"
)

type ICarRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, input *entity.Car) (*entity.Car, error)
	GetAll(ctx context.Context, tx *sql.Tx) ([]entity.Car, error)
	GetDetail(ctx context.Context, tx *sql.Tx, id int) (*entity.Car, error)
}
