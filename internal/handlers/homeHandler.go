package handlers

import (
	"cars/internal/api"
	"cars/internal/pagination"
	"html/template"
	"log"
	"net/http"
)

const baseURL = "http://localhost:3000/api"

func HomeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := pagination.PageData{}

		errCh := make(chan error, 3)
		doneCh := make(chan struct{}, 3)

		go func() {
			err := api.FetchJSON(baseURL+"/models", &data.Cars)
			errCh <- err
			doneCh <- struct{}{}
		}()

		go func() {
			err := api.FetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
			errCh <- err
			doneCh <- struct{}{}
		}()

		go func() {
			err := api.FetchJSON(baseURL+"/categories", &data.Categories)
			errCh <- err
			doneCh <- struct{}{}
		}()

		for i := 0; i < 3; i++ {
			<-doneCh
		}

		close(errCh)

		for err := range errCh {
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		data.CompareIDs = GetCompareIDMap(r)
		data.CompareCount = len(GetCompareIDs(r))
		data.MaxCompareCars = MaxCompareCars

		pagination.PaginateCars(&data, r)

		err := tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
