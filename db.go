package main

import (
	"errors"
	"fmt"

	"gitlab.com/my-repos/go-seed-db/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	ErrNoDatabaseConfig = errors.New(`no "db-addr" value in config`)
	ErrNoDatabase       = errors.New("no database in API")
	ErrDatabaseFound    = errors.New("unexpected database in API")
)

func (api *API) ConnectDB(dbaddr string) (func(), error) {
	if dbaddr == "" {
		return nil, ErrNoDatabaseConfig
	}

	api.logger.Println("Connecting to Database")

	db, err := gorm.Open(postgres.Open(dbaddr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: ".", SingularTable: false}})
	if err != nil {
		return nil, fmt.Errorf("unable to open Database: %w", err)
	}

	api.db = db

	//Closer function for cleaning up DB cnnection
	return func() {
		sqlDB, err := api.db.DB()
		if err != nil {
			api.logger.Printf("\nerror getting sql database: %v", err)
			return
		}

		sqlDB.Close()
		api.db = nil
		api.logger.Println("Disconnected from Database")
	}, nil
}

// ListDBTrainers lists all trainers from postgres DB
// GORM pluralizes struct name to snake_cases as table name, for eg, struct User, table name is users by convention
// GORM Column db name uses the field’s name’s snake_case by convention. for eg. CreatedAt is created_at by convention if not explicitly mention in model tags

func (api *API) ListDBTrainers() ([]*model.Trainer, error) {
	trainers := []*model.Trainer{}
	query := api.db.Table("trainers") // dont need to do this
	if err := query.Find(&trainers).Error; err != nil {
		return nil, err
	}

	return trainers, nil
}

func (api *API) GetDBTrainerByName(name string) (*model.Trainer, error) {
	trainer := &model.Trainer{}

	query := api.db.Table("trainers").Where("name=?", name)

	if err := query.First(trainer).Error; err != nil {
		return nil, err
	}

	return trainer, nil
}

func (api *API) GetDBTrainerByID(id string) (*model.Trainer, error) {
	trainer := &model.Trainer{}

	query := api.db.Table("trainers") // dont need to do this

	if err := query.First(trainer, id).Error; err != nil {
		return nil, err
	}

	return trainer, nil
}

func (api *API) CreateDBTrainer(t *model.Trainer) error {
	query := api.db.Table("trainers") // dont need to do this

	if err := query.Create(t).Error; err != nil {
		return err
	}

	return nil
}

func (api *API) UpdateDBTrainer(t *model.Trainer) error {
	query := api.db.Table("trainers") // dont need to do this

	tx := query.Model(t).Where("id = ?", t.ID).Updates(t)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("No trainers with id:%s", t.ID)
	}

	return nil
}

func (api *API) DeleteDBTrainer(id string) error {
	query := api.db.Table("trainers")

	tx := query.Delete(&model.Trainer{ID: id})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("No trainers with id:%s", id)
	}

	return nil
}
