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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
	Pages      []PageLink
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{}

	errCh := make(chan error, 3) // буфер важен!
	doneCh := make(chan struct{}, 3)

	go func() {
		err := fetchJSON(baseURL+"/models", &data.Cars)
		errCh <- err
		doneCh <- struct{}{}
	}()

	go func() {
		err := fetchJSON(baseURL+"/manufacturers", &data.Manufacturers)
		errCh <- err
		doneCh <- struct{}{}
	}()

	go func() {
		err := fetchJSON(baseURL+"/categories", &data.Categories)
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

	paginateCars(&data, r)

	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
	Page         int
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

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	specPage := CarWithDetails{
		Car:          car,
		Manufacturer: manufacturer,
		Category:     category,
		Page:         page,
	}

	err := tmpl.ExecuteTemplate(w, "car.html", specPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PageLink struct {
	Number int
	IsDots bool
}

func buildPages(currentPage, totalPages int) []PageLink {
	const sidePages = 2

	var pages []PageLink

	addPage := func(n int) {
		pages = append(pages, PageLink{
			Number: n,
			IsDots: false,
		})
	}

	addDots := func() {
		pages = append(pages, PageLink{
			IsDots: true,
		})
	}

	if totalPages <= 1 {
		addPage(1)
		return pages
	}

	// Если страниц мало, показываем все
	if totalPages <= 7 {
		for i := 1; i <= totalPages; i++ {
			addPage(i)
		}
		return pages
	}

	addPage(1)

	start := currentPage - sidePages
	end := currentPage + sidePages

	if start < 2 {
		start = 2
	}

	if end > totalPages-1 {
		end = totalPages - 1
	}

	if start > 2 {
		addDots()
	}

	for i := start; i <= end; i++ {
		addPage(i)
	}

	if end < totalPages-1 {
		addDots()
	}

	addPage(totalPages)

	return pages
}

func paginateCars(data *PageData, r *http.Request) {
	const perPage = 20

	page := getCurrentPage(r)

	totalCars := len(data.Cars)
	totalPages := calculateTotalPages(totalCars, perPage)

	if page > totalPages {
		page = totalPages
	}

	start, end := getPageBounds(page, perPage, totalCars)

	data.Cars = data.Cars[start:end]

	data.Page = page
	data.TotalPages = totalPages
	data.PrevPage = page - 1
	data.NextPage = page + 1
	data.HasPrev = page > 1
	data.HasNext = page < totalPages
	data.Pages = buildPages(page, totalPages)
}

func getCurrentPage(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}

	return page
}

func calculateTotalPages(totalItems, perPage int) int {
	if totalItems == 0 {
		return 1
	}

	return int(math.Ceil(float64(totalItems) / float64(perPage)))
}

func getPageBounds(page, perPage, totalItems int) (int, int) {
	start := (page - 1) * perPage
	end := start + perPage

	if end > totalItems {
		end = totalItems
	}

	return start, end
}
