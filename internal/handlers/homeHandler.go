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

		go func() {
			errCh <- api.FetchJSON(baseURL+"/models", &data.Cars)
		}()

		go func() {
			errCh <- api.FetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
		}()

		go func() {
			errCh <- api.FetchJSON(baseURL+"/categories", &data.Categories)
		}()

		for i := 0; i < 3; i++ {
			if err := <-errCh; err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				tmpl.ExecuteTemplate(w, "maintenance.html", nil)
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
