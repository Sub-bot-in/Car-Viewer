package handlers

import (
	"cars/catalog"
	"cars/internal/api"
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
)

type CarWithDetails struct {
	Car          catalog.Model
	Manufacturer catalog.Manufacturers
	Category     catalog.Categories
	RelatedCars  []catalog.Model
	Page         int

	ManufacturerFilter string
	CategoryFilter     string
	YearFilter         string

	IsCompared     bool
	CompareCount   int
	MaxCompareCars int
}

func SpecHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/specifications/")

		data := api.LoadCarsAPI()

		var car catalog.Model
		var manufacturer catalog.Manufacturers
		var category catalog.Categories

		for _, c := range data.Cars {
			if fmt.Sprint(c.ID) == id {
				car = c
				break
			}
		}

		for _, c := range data.Manufacturers {
			if c.ID == car.ManufacturerID {
				manufacturer = c
				break
			}
		}

		for _, c := range data.Categories {
			if c.ID == car.CategoryID {
				category = c
				break
			}
		}

		page := 1

		pageStr := r.URL.Query().Get("page")
		if pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err == nil && p > 0 {
				page = p
			}
		}

		relatedCars := getRelatedCars(data.Cars, car, 10)

		compareIDs := GetCompareIDMap(r)
		compareList := GetCompareIDs(r)

		specPage := CarWithDetails{
			Car:          car,
			Manufacturer: manufacturer,
			Category:     category,
			RelatedCars:  relatedCars,
			Page:         page,

			ManufacturerFilter: r.URL.Query().Get("manufacturer"),
			CategoryFilter:     r.URL.Query().Get("category"),
			YearFilter:         r.URL.Query().Get("year"),

			IsCompared:     compareIDs[car.ID],
			CompareCount:   len(compareList),
			MaxCompareCars: MaxCompareCars,
		}

		err := tmpl.ExecuteTemplate(w, "car.html", specPage)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getRelatedCars(cars []catalog.Model, currentCar catalog.Model, limit int) []catalog.Model {
	var related []catalog.Model

	for _, car := range cars {
		if car.ID == currentCar.ID {
			continue
		}

		if car.ManufacturerID == currentCar.ManufacturerID && car.ID != currentCar.ID {
			related = append(related, car)
		}
	}

	rand.Shuffle(len(related), func(i, j int) {
		related[i], related[j] = related[j], related[i]
	})

	if len(related) > limit {
		related = related[:limit]
	}

	return related
}
