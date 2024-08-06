package db

import (
	"context"

	"github.com/SevvyP/items/pkg"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateItem(item pkg.Item, context context.Context) error {
	args := m.Called(item, context)
	return args.Error(0)
}

func (m *MockDB) GetItem(name string, context context.Context) (*pkg.Item, error) {
	args := m.Called(name, context)
	return args.Get(0).(*pkg.Item), args.Error(1)
}

func (m *MockDB) UpdateItem(item pkg.Item, context context.Context) error {
	args := m.Called(item, context)
	return args.Error(0)
}

func (m *MockDB)DeleteItem(name string, context context.Context) (*pkg.Item, error) {
	args := m.Called(name, context)
	return args.Get(0).(*pkg.Item), args.Error(1)
}