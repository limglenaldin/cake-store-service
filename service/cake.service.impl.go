package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/glenaldinlim/cake-store-service/exception"
	"github.com/glenaldinlim/cake-store-service/model/entity"
	"github.com/glenaldinlim/cake-store-service/model/web"
	"github.com/glenaldinlim/cake-store-service/repository"
	"github.com/glenaldinlim/cake-store-service/utils"
	"github.com/go-playground/validator/v10"
)

type CakeServiceImpl struct {
	CakeRepository repository.CakeRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewCakeService(cakeRepository repository.CakeRepository, DB *sql.DB, validate *validator.Validate) CakeService {
	return &CakeServiceImpl{
		CakeRepository: cakeRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (service *CakeServiceImpl) FindAll(ctx context.Context) []entity.Cake {
	tx, err := service.DB.Begin()
	if err != nil {
		utils.Logger().Errorf("[DB Tx Exception] CakeService.FindAll: %s", err.Error())
		panic(err)
	}
	defer utils.CommitOrRollback(tx)

	cakes := service.CakeRepository.FindAll(ctx, tx)
	return cakes
}

func (service *CakeServiceImpl) FindById(ctx context.Context, cakeId int64) entity.Cake {
	tx, err := service.DB.Begin()
	if err != nil {
		utils.Logger().Errorf("[DB Tx Exception] CakeService.FindById: %s", err.Error())
		panic(err)
	}
	defer utils.CommitOrRollback(tx)

	cake, err := service.CakeRepository.FindById(ctx, tx, cakeId)
	if err != nil {
		panic(exception.NewNotFounderror(err.Error()))
	}

	return cake
}

func (service *CakeServiceImpl) Create(ctx context.Context, request web.CakeRequest) entity.Cake {
	err := service.Validate.Struct(request)
	if err != nil {
		utils.Logger().Errorf("[Validation] CakeService.Create: %s", err.Error())
		panic(err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		utils.Logger().Errorf("[DB Tx Exception] CakeService.Create: %s", err.Error())
		panic(err)
	}
	defer utils.CommitOrRollback(tx)

	cake := entity.Cake{
		Title:       request.Title,
		Description: request.Description,
		Rating:      request.Rating,
		Image:       request.Image,
	}

	cake = service.CakeRepository.Save(ctx, tx, cake)
	return cake
}

func (service *CakeServiceImpl) Update(ctx context.Context, request web.CakeRequest, cakeId int64) entity.Cake {
	err := service.Validate.Struct(request)
	if err != nil {
		utils.Logger().Errorf("[Validation] CakeService.Update: %s", err.Error())
		panic(err)
	}

	tx, err := service.DB.Begin()
	if err != nil {
		utils.Logger().Errorf("[DB Tx Exception] CakeService.Update: %s", err.Error())
		panic(err)
	}
	defer utils.CommitOrRollback(tx)

	cake, err := service.CakeRepository.FindById(ctx, tx, cakeId)
	if err != nil {
		panic(exception.NewNotFounderror(err.Error()))
	}

	cake.Title = request.Title
	cake.Description = request.Description
	cake.Rating = request.Rating
	cake.Image = request.Image
	cake.UpdatedAt = time.Now()

	cake = service.CakeRepository.Update(ctx, tx, cake)

	return cake
}

func (service *CakeServiceImpl) Delete(ctx context.Context, cakeId int64) {
	tx, err := service.DB.Begin()
	if err != nil {
		utils.Logger().Errorf("[DB Tx Exception] CakeService.Delete: %s", err.Error())
		panic(err)
	}
	defer utils.CommitOrRollback(tx)

	cake, err := service.CakeRepository.FindById(ctx, tx, cakeId)
	if err != nil {
		panic(exception.NewNotFounderror(err.Error()))
	}

	service.CakeRepository.Delete(ctx, tx, cake)
}
