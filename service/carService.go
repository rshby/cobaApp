package service

import (
	"cobaApp/customError"
	"cobaApp/helper"
	"cobaApp/model/dto"
	"cobaApp/model/entity"
	"cobaApp/repository"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type CarService struct {
	DB            *sql.DB
	Validate      *validator.Validate
	CarRepository repository.ICarRepository
}

// function provider
func NewCarService(db *sql.DB, validate *validator.Validate, carRepo repository.ICarRepository) ICarService {
	return &CarService{
		DB:            db,
		Validate:      validate,
		CarRepository: carRepo,
	}
}

func (c *CarService) Insert(ctx context.Context, request *dto.InsertCarRequest) (*dto.InsertCarResponse, error) {
	// start tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Service Insert")
	defer span.Finish()

	reqJson, _ := json.Marshal(&request)
	span.LogFields(log.String("request", string(reqJson)))

	if err := c.Validate.StructCtx(ctxTracing, *request); err != nil {
		// return error validator
		return nil, err
	}

	// create entity input
	input := entity.Car{
		Name:  request.Name,
		Price: request.Price,
	}

	// validate date if null
	if request.ReleaseDate == "" {
		// jika release_date nya null
		input.ReleaseDate = &sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	} else {
		// jika release_date diisi
		input.ReleaseDate = &sql.NullTime{
			Time:  helper.StringToDate(request.ReleaseDate),
			Valid: true,
		}
	}

	// create db transaction
	tx, _ := c.DB.Begin()
	defer tx.Rollback()

	// run query in repository
	result, err := c.CarRepository.Insert(ctxTracing, tx, &input)
	if err != nil {
		return nil, err
	}

	// success insert
	tx.Commit()

	// create respone
	response := dto.InsertCarResponse{
		Id:    result.Id,
		Name:  result.Name,
		Price: result.Price,
	}

	if result.ReleaseDate.Valid {
		response.ReleaseDate = result.ReleaseDate.Time.Format("2006-01-02")
	} else {
		response.ReleaseDate = ""
	}

	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))

	// return response
	return &response, nil
}

func (c *CarService) GetAll(ctx context.Context) ([]dto.InsertCarResponse, error) {
	// start tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Service GetAll")
	defer span.Finish()

	// create transaction
	tx, _ := c.DB.Begin()
	defer tx.Rollback()

	// run query in repository
	cars, err := c.CarRepository.GetAll(ctxTracing, tx)
	if err != nil {
		return nil, err
	}

	// commit transaction
	tx.Commit()

	// convert to response
	var response = []dto.InsertCarResponse{}
	for _, data := range cars {
		response = append(response, dto.InsertCarResponse{
			Id:          data.Id,
			Name:        data.Name,
			Price:       data.Price,
			ReleaseDate: helper.DateToString(data.ReleaseDate.Time),
		})
	}

	// log to tracing
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))

	// return all response
	return response, nil
}

func (c *CarService) GetDetail(ctx context.Context, id int) (*dto.InsertCarResponse, error) {
	// create span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Service GetDetail")
	defer span.Finish()

	span.LogFields(log.Int("id", id))

	// start transaction
	tx, err := c.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		return nil, customError.NewInternalServerError(err.Error())
	}

	// call procedure in repository
	car, err := c.CarRepository.GetDetail(ctxTracing, tx, id)
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		return nil, err
	}

	// convert to response dto
	response := dto.InsertCarResponse{
		Id:          car.Id,
		Name:        car.Name,
		Price:       car.Price,
		ReleaseDate: helper.DateToString(car.ReleaseDate.Time),
	}

	// log to tracing
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))

	tx.Commit()
	return &response, nil
}
