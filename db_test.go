package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sshetty10/go-seed-db/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var id string
var resolver *Resolver
var api *API

func TestMain(m *testing.M) {
	log.Println("Stuff BEFORE the tests!")

	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("Stuff AFTER the tests!")
	os.Exit(code)

}

func run(m *testing.M) (code int, err error) {
	// pseudo-code, some implementation excluded:
	//
	db, dbCloser, _, err := InitTestDatabase("mydb")
	if err != nil {
		return -1, err
	}

	if dbCloser != nil {
		defer dbCloser()
	}

	logger := log.New(os.Stdout, "go-seed-test", log.Ldate|log.Ltime|log.Llongfile)

	api = &API{db: db, logger: logger}

	resolver = api.NewResolver()

	return m.Run(), nil
}

func TestListDBTrainers(t *testing.T) {
	trainers, err := api.ListDBTrainers()

	require.NoError(t, err)
	assert.Equal(t, len(trainers), 4)
}

func TestCreateDBTrainer(t *testing.T) {
	tr := &model.Trainer{
		Name:      "sometest",
		City:      "somecity",
		Age:       70,
		LicenseID: "TX-38274",
	}
	err := api.CreateDBTrainer(tr)
	require.NoError(t, err)
	assert.NotEqual(t, tr.ID, "")
	id = tr.ID
}

func TestUpdateDBTrainer(t *testing.T) {
	tr := &model.Trainer{
		ID:        id,
		Name:      "somenewtest",
		City:      "somecity",
		Age:       70,
		LicenseID: "TX-38274",
	}
	err := api.UpdateDBTrainer(tr)
	require.NoError(t, err)
	assert.Equal(t, tr.Name, "somenewtest")
}

func TestGetDBTrainerByID(t *testing.T) {
	var newtr *model.Trainer
	tr := &model.Trainer{
		ID:        id,
		Name:      "somenewtest",
		City:      "somecity",
		Age:       70,
		LicenseID: "TX-38274",
	}
	newtr, err := api.GetDBTrainerByID(id)
	require.NoError(t, err)
	assert.Equal(t, tr.ID, newtr.ID)
	assert.Equal(t, tr.Name, newtr.Name)
}

func TestGetDBTrainerByName(t *testing.T) {
	var newtr *model.Trainer
	tr := &model.Trainer{
		ID:        id,
		Name:      "somenewtest",
		City:      "somecity",
		Age:       70,
		LicenseID: "VA-38274",
	}
	newtr, err := api.GetDBTrainerByName("somenewtest")
	require.NoError(t, err)
	assert.Equal(t, tr.ID, newtr.ID)
	assert.Equal(t, tr.Name, newtr.Name)
}
