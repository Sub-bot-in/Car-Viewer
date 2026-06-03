package main

import (
	"cars/catalog"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://localhost:3000/api"

func fetchJSON(url string, target interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("bad status: %d", resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		panic(err)
	}
}

func fetchModels() []catalog.Model {
	var models []catalog.Model
	fetchJSON(baseURL+"/models", &models)
	return models
}

func fetchManufacturers() []catalog.Manufacturers {
	var manufacturers []catalog.Manufacturers
	fetchJSON(baseURL+"/manufacturers", &manufacturers)
	return manufacturers
}

func fetchCategories() []catalog.Categories {
	var categories []catalog.Categories
	fetchJSON(baseURL+"/categories", &categories)
	return categories
}
