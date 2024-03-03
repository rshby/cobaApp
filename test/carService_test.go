package test

import (
	"cobaApp/customError"
	"cobaApp/helper"
	"cobaApp/model/dto"
	"cobaApp/model/entity"
	"cobaApp/service"
	mck "cobaApp/test/mock"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var validate = validator.New()

func TestInsertCar(t *testing.T) {
	t.Run("test insert error cant insert", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		defer db.Close()
		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()
		errMessage := "error when add new data"
		carRepo.Mock.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(nil, customError.NewInternalServerError(errMessage))

		// test
		result, err := carService.Insert(context.Background(), &dto.InsertCarRequest{
			Name:        "Toyota",
			Price:       416000000,
			ReleaseDate: "2022-10-10",
		})

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
	})
	t.Run("test insert new car success", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		defer db.Close()
		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()
		carRepo.Mock.On("Insert", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Car{
			Id:    1,
			Name:  "Toyota",
			Price: 500000000,
			ReleaseDate: &sql.NullTime{
				Time:  helper.StringToDate("2020-10-10"),
				Valid: true,
			},
		}, nil)

		// test
		result, err := carService.Insert(context.Background(), &dto.InsertCarRequest{
			Name:        "Toyota",
			Price:       500000000,
			ReleaseDate: "2020-10-10",
		})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Id)
		assert.Equal(t, "Toyota", result.Name)
		carRepo.Mock.AssertExpectations(t)
	})
	t.Run("test insert new car date empty", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// test
		result, err := carService.Insert(context.Background(), &dto.InsertCarRequest{
			Name:        "Toyota",
			Price:       1,
			ReleaseDate: "",
		})

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		carRepo.Mock.AssertExpectations(t)
	})
}

func TestGetAllCar(t *testing.T) {
	t.Run("test get all not found", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		defer db.Close()

		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()
		errMessage := "record not found"
		carRepo.Mock.On("GetAll", mock.Anything, mock.Anything).Return(nil, customError.NewNotFoundError(errMessage))

		// test
		cars, err := carService.GetAll(context.Background())
		assert.Nil(t, cars)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
		carRepo.Mock.AssertExpectations(t)
	})
	t.Run("test get all cars success", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		defer db.Close()

		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()
		response := []entity.Car{
			{
				Id:    1,
				Name:  "Toyota",
				Price: 1,
				ReleaseDate: &sql.NullTime{
					Time:  helper.StringToDate("2020-10-10"),
					Valid: true,
				},
			},
			{
				Id:    2,
				Name:  "Honda",
				Price: 1,
				ReleaseDate: &sql.NullTime{
					Time:  helper.StringToDate("2020-10-10"),
					Valid: true,
				},
			},
		}
		carRepo.Mock.On("GetAll", mock.Anything, mock.Anything).
			Return(response, nil)

		// test
		cars, err := carService.GetAll(context.Background())

		assert.Nil(t, err)
		assert.NotNil(t, cars)
		assert.Equal(t, 2, len(cars))
	})
}

func TestGetDetailCar(t *testing.T) {
	t.Run("test get detail not found", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		defer db.Close()

		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		errMessage := "record not found"
		carRepo.Mock.On("GetDetail", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, customError.NewNotFoundError(errMessage))

		// test
		car, err := carService.GetDetail(context.Background(), 1)

		assert.Nil(t, car)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errMessage, err.Error())
	})
	t.Run("test get detail success", func(t *testing.T) {
		db, dbMock, _ := sqlmock.New()
		defer db.Close()

		carRepo := mck.NewCarRepositoryMock()
		carService := service.NewCarService(db, validate, carRepo)

		// mock
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		carRepo.Mock.On("GetDetail", mock.Anything, mock.Anything, mock.Anything).
			Return(&entity.Car{
				Id:    1,
				Name:  "Toyota",
				Price: 450000000,
				ReleaseDate: &sql.NullTime{
					Time:  helper.StringToDate("2020-10-10"),
					Valid: true,
				},
			}, nil)

		// test
		car, err := carService.GetDetail(context.Background(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, car)
		assert.Equal(t, 1, car.Id)
		assert.Equal(t, "Toyota", car.Name)
	})
}
