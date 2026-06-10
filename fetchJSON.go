package main

import (
	"cars/catalog"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://localhost:3000/api"

func fetchJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}
}

modelsCh := make(chan []catalog.Model)
//brandsCh := make(chan []catalog.Brand)
errCh := make(chan error)

go fetchModels(modelsCh, errCh)

func fetchModels(modelsCh, errCh) []catalog.Model {
	var models []catalog.Model
	if err := fetchJSON(baseURL+"/models", &models); err != nil {
		errCh <- err
		return
	}
	modelsCh <- models
}

func fetchManufacturers(manufacturersCh, errCh) []catalog.Manufacturers {
	var manufacturers []catalog.Manufacturers
	fetchJSON(baseURL+"/manufacturers", &manufacturers)
	return manufacturers
}

func fetchCategories(categoriesCh, errCh) []catalog.Categories {
	var categories []catalog.Categories
	fetchJSON(baseURL+"/categories", &categories)
	return categories
}
