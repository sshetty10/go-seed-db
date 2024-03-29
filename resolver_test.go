package main

import (
	"context"
	"testing"

	"github.com/sshetty10/go-seed-db/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrainers(t *testing.T) {
	ctx := context.Background()
	trainers, err := resolver.Query().Trainers(ctx)

	require.NoError(t, err)
	assert.Equal(t, len(trainers), 5)
	assert.Equal(t, trainers[0].LicenseState, "VA")
}

func TestLicenseState(t *testing.T) {
	ctx := context.Background()
	tr := &model.Trainer{
		ID:        id,
		Name:      "somenewtest",
		City:      "somecity",
		Age:       70,
		LicenseID: "VA-38274",
	}
	st, err := resolver.Trainer().LicenseState(ctx, tr)

	require.NoError(t, err)
	assert.Equal(t, st, "VA")
}

func BenchmarkLicenseState(b *testing.B) {
	ctx := context.Background()
	tr := &model.Trainer{
		ID:        id,
		Name:      "somenewtest",
		City:      "somecity",
		Age:       70,
		LicenseID: "VA-38274",
	}
	for i := 0; i < b.N; i++ {
		resolver.Trainer().LicenseState(ctx, tr)
	}
}

/*
// Benchmarking with a range of inputs

func benchmarkLicenseState(tr *model.Trainer, b *testing.B) {
    ctx := context.Background()
    for i := 0; i < b.N; i++ {
        resolver.Trainer().LicenseState(ctx, tr)
    }
}

func BenchmarkLicenseState1(b *testing.B)  { benchmarkLicenseState(&model.Trainer{ID:id,Name:"foo",City:"bar",Age:70,LicenseID: "VA-38274",}, b) }
func BenchmarkLicenseState2(b *testing.B)  { benchmarkLicenseState(&model.Trainer{ID:id,Name:"foo",City:"bar",Age:70,LicenseID: "MD-38274",}, b) }
func BenchmarkLicenseState3(b *testing.B)  { benchmarkLicenseState(&model.Trainer{ID:id,Name:"foo",City:"bar",Age:70,LicenseID: "NY-38274",}, b) }
*/
