package db

import (
	"errors"

	"github.com/SevvyP/plants/pkg"
)

type DB struct {
	plants []*pkg.Plant
}

func NewDB() *DB {
	return &DB{[]*pkg.Plant{{Name: "Cactus", Description: "pointy green thing"}}}
}

func (db *DB) GetPlant(name string) (*pkg.Plant, error){
	for _, plant := range db.plants {
		if plant.Name == name {
			return plant, nil
		}
	}
	return nil, errors.New("plant not found")
}