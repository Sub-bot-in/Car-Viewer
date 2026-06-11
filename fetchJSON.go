package main

import (
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
	return nil
}

func fetch[T any](url string, resultCh chan<- T, errCh chan<- error) {
	var data T

	if err := fetchJSON(url, &data); err != nil {
		errCh <- err
		return
	}

	resultCh <- data
}
