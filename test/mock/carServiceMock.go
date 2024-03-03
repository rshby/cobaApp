package mock

import (
	"cobaApp/model/dto"
	"context"
	"github.com/stretchr/testify/mock"
)

type CarServiceMock struct {
	Mock mock.Mock
}

// function provider
func NewCarServiceMock() *CarServiceMock {
	return &CarServiceMock{mock.Mock{}}
}

func (c *CarServiceMock) Insert(ctx context.Context, request *dto.InsertCarRequest) (*dto.InsertCarResponse, error) {
	args := c.Mock.Called(ctx, request)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*dto.InsertCarResponse), nil
}

func (c *CarServiceMock) GetAll(ctx context.Context) ([]dto.InsertCarResponse, error) {
	args := c.Mock.Called(ctx)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.([]dto.InsertCarResponse), nil
}

func (c *CarServiceMock) GetDetail(ctx context.Context, id int) (*dto.InsertCarResponse, error) {
	args := c.Mock.Called(ctx, id)

	value := args.Get(0)
	if value == nil {
		return nil, args.Error(1)
	}

	return value.(*dto.InsertCarResponse), nil
}
