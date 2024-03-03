package service

import (
	"cobaApp/model/dto"
	"context"
)

type ICarService interface {
	Insert(ctx context.Context, request *dto.InsertCarRequest) (*dto.InsertCarResponse, error)
	GetAll(ctx context.Context) ([]dto.InsertCarResponse, error)
	GetDetail(ctx context.Context, id int) (*dto.InsertCarResponse, error)
}
