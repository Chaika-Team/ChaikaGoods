package main

import (
	"ChaikaGoods/internal/models"
	"encoding/json"
	"testing"
)

func TestProductJSON(m *testing.T) {
	// test models.Product to json
	p := models.Product{
		ID:          new(int64),
		Name:        new(string),
		Description: new(string),
		Price:       new(float64),
		ImageURL:    new(string),
		SKU:         new(string),
	}
	// set name
	*p.Name = "test"

	// convert to json
	j, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	println(string(j))
}
