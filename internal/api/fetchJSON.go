package api

import (
	"cars/catalog"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const baseURL = "http://localhost:3000/api"

func FetchJSON(url string, target interface{}) error {
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
	return nil
}

func fetch[T any](url string, resultCh chan<- T, errCh chan<- error) {
	var data T

	if err := FetchJSON(url, &data); err != nil {
		errCh <- err
		return
	}

	resultCh <- data
}

type PageData struct {
	Cars          []catalog.Model
	Manufacturers []catalog.Manufacturers
	Categories    []catalog.Categories
}

func LoadCarsAPI() PageData {
	var data PageData

	err := FetchJSON(baseURL+"/models", &data.Cars)
	if err != nil {
		log.Println(err)
	}

	err = FetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
	if err != nil {
		log.Println(err)
	}

	err = FetchJSON(baseURL+"/categories", &data.Categories)
	if err != nil {
		log.Println(err)
	}

	return data
}
