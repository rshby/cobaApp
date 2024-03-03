package mock

import (
	"cobaApp/model/entity"
	"context"
	"database/sql"
	"github.com/stretchr/testify/mock"
)

type CarRepositoryMock struct {
	Mock mock.Mock
}

// function provider
func NewCarRepositoryMock() *CarRepositoryMock {
	return &CarRepositoryMock{Mock: mock.Mock{}}
}

func (c *CarRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, input *entity.Car) (*entity.Car, error) {
	arguments := c.Mock.Called(ctx, tx)

	value := arguments.Get(0)
	if value == nil {
		return nil, arguments.Error(1)
	}

	return value.(*entity.Car), nil
}

func (c *CarRepositoryMock) GetAll(ctx context.Context, tx *sql.Tx) ([]entity.Car, error) {
	args := c.Mock.Called(ctx, tx)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.([]entity.Car), nil
}

func (c *CarRepositoryMock) GetDetail(ctx context.Context, tx *sql.Tx, id int) (*entity.Car, error) {
	args := c.Mock.Called(ctx, id)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*entity.Car), nil
}
