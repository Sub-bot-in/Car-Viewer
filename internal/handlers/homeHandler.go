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

		ch := make(chan error, 3)

		go func() {
			ch <- api.FetchJSON(baseURL+"/models", &data.Cars)
		}()

		go func() {
			ch <- api.FetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
		}()

		go func() {
			ch <- api.FetchJSON(baseURL+"/categories", &data.Categories)
		}()

		for i := 0; i < 3; i++ { //program waits when all three API calls are completed
			if err := <-ch; err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)

				if err := tmpl.ExecuteTemplate(w, "maintenance.html", nil); err != nil {
					log.Println(err)
				}

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
