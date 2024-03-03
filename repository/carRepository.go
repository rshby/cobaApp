package repository

import (
	"cobaApp/customError"
	"cobaApp/model/entity"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CarRepository struct {
	DB *sql.DB
}

// function provider
func NewCarRepository(db *sql.DB) ICarRepository {
	return &CarRepository{
		DB: db,
	}
}

// method implementasi Insert
func (c *CarRepository) Insert(ctx context.Context, tx *sql.Tx, input *entity.Car) (*entity.Car, error) {
	// start tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Repository Insert")
	defer span.Finish()

	reqString, _ := json.Marshal(&input)
	span.LogFields(log.String("request", string(reqString)))

	// prepare query
	statement, err := tx.PrepareContext(ctxTracing, "INSERT INTO cars(name, price, release_date) VALUES (?, ?, ?)")
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	result, err := statement.ExecContext(ctxTracing, input.Name, input.Price, input.ReleaseDate.Time)
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	id, err := result.RowsAffected()
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	// success insert
	input.Id = int(id)
	return input, nil
}

// method implementasi GetAll
func (c *CarRepository) GetAll(ctx context.Context, tx *sql.Tx) ([]entity.Car, error) {
	// start tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Repository GetAll")
	defer span.Finish()

	// prepare query
	statement, err := tx.PrepareContext(ctxTracing, "SELECT id, name, price, release_date FROM cars")
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	// execute query
	rows, err := statement.QueryContext(ctxTracing)
	if err != nil {
		return nil, customError.NewInternalServerError(err.Error())
	}

	var response []entity.Car
	for rows.Next() {
		var res entity.Car
		if err := rows.Scan(&res.Id, &res.Name, &res.Price, &res.ReleaseDate); err != nil {
			if err == sql.ErrNoRows {
				return nil, customError.NewNotFoundError("record not found")
			}

			return nil, customError.NewInternalServerError(err.Error())
		}

		response = append(response, res)
	}

	// if not found
	if len(response) == 0 {
		return nil, customError.NewNotFoundError("record not found")
	}

	// log response to tracing
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))

	// success get data
	return response, nil
}

// method implementasi get detail by id
func (c *CarRepository) GetDetail(ctx context.Context, tx *sql.Tx, id int) (*entity.Car, error) {
	// start span tracing
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Repository GetDetail")
	defer span.Finish()

	// log id to tracing
	span.LogFields(log.Int("id", id))

	// prepare query
	statement, err := tx.PrepareContext(ctxTracing, "SELECT id, name, price, release_date FROM cars WHERE id=?")
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		return nil, customError.NewInternalServerError(err.Error())
	}

	// query
	row := statement.QueryRowContext(ctxTracing, id)
	if row.Err() != nil {
		span.LogFields(log.String("error", row.Err().Error()))

		if row.Err() == sql.ErrNoRows {
			return nil, customError.NewNotFoundError(row.Err().Error())
		}

		return nil, customError.NewInternalServerError(row.Err().Error())
	}

	var response entity.Car
	if err := row.Scan(&response.Id, &response.Name, &response.Price, &response.ReleaseDate); err != nil {
		span.LogFields(log.String("error", err.Error()))
		if err == sql.ErrNoRows {
			return nil, customError.NewNotFoundError(err.Error())
		}

		return nil, customError.NewInternalServerError(err.Error())
	}

	// success get data
	resJson, _ := json.Marshal(&response)
	span.LogFields(log.String("response", string(resJson)))
	return &response, nil
}
