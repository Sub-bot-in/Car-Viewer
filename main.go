package main

import (
	"cars/catalog"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const baseURL = "http://localhost:3000/api"

func main() {

	http.HandleFunc("/", homeHandler)
	//	http.HandleFunc("/cars", carsHandler)
	http.HandleFunc("/specifications/", specHandler)
	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var tmpl *template.Template = template.Must(template.ParseGlob("templates/*.html"))

type PageData struct {
	Cars          []catalog.Model
	Manufacturers []catalog.Manufacturers
	Categories    []catalog.Categories

	Page       int
	TotalPages int
	PrevPage   int
	NextPage   int
	HasPrev    bool
	HasNext    bool
	Pages      []int
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{}

	errCh := make(chan error, 3)

	go func() {
		errCh <- fetchJSON(baseURL+"/models", &data.Cars)
	}()

	go func() {
		errCh <- fetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
	}()

	go func() {
		errCh <- fetchJSON(baseURL+"/categories", &data.Categories)
	}()
<<<<<<< HEAD

=======
>>>>>>> 3c5de55 (api fixed)
	for range 3 {
		err := <-errCh
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	const perPage = 20

	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	totalCars := len(data.Cars)
	totalPages := int(math.Ceil(float64(totalCars) / float64(perPage)))

	if totalPages == 0 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	end := start + perPage

	if end > totalCars {
		end = totalCars
	}

	data.Cars = data.Cars[start:end]

	data.Page = page
	data.TotalPages = totalPages
	data.PrevPage = page - 1
	data.NextPage = page + 1
	data.HasPrev = page > 1
	data.HasNext = page < totalPages

	for i := 1; i <= totalPages; i++ {
		data.Pages = append(data.Pages, i)
	}

	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func carsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("cars"))
// }

func loadCarsAPI() PageData {
	var data PageData

	err := fetchJSON(baseURL+"/models", &data.Cars)
	if err != nil {
		log.Println(err)
	}

	err = fetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
	if err != nil {
		log.Println(err)
	}

	err = fetchJSON(baseURL+"/categories", &data.Categories)
	if err != nil {
		log.Println(err)
	}

	return data
}

type CarWithDetails struct {
	Car          catalog.Model
	Manufacturer catalog.Manufacturers
	Category     catalog.Categories
}

func specHandler(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/specifications/")

	data := loadCarsAPI()

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
		if fmt.Sprint(c.ID) == id {
			manufacturer = c
			break
		}
	}

	for _, c := range data.Categories {
		if fmt.Sprint(c.ID) == id {
			category = c
			break
		}
	}

	specPage := CarWithDetails{
		Car:          car,
		Manufacturer: manufacturer,
		Category:     category,
	}
	err := tmpl.ExecuteTemplate(w, "car.html", specPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
