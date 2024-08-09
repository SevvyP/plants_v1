package db

import (
	"context"

	"github.com/SevvyP/plants/pkg"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreatePlant(plant pkg.Plant, context context.Context) error {
	args := m.Called(plant, context)
	return args.Error(0)
}

func (m *MockDB) GetPlant(name string, context context.Context) (*pkg.Plant, error) {
	args := m.Called(name, context)
	return args.Get(0).(*pkg.Plant), args.Error(1)
}

func (m *MockDB) UpdatePlant(plant pkg.Plant, context context.Context) error {
	args := m.Called(plant, context)
	return args.Error(0)
}

func (m *MockDB)DeletePlant(name string, context context.Context) (*pkg.Plant, error) {
	args := m.Called(name, context)
	return args.Get(0).(*pkg.Plant), args.Error(1)
}